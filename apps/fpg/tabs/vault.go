// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package tabs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/grd/FreePDM/internal/vaults"
)

var reNumeric = regexp.MustCompile(`^\d+$`)

type VaultTab struct {
	Tab       *container.TabItem
	Root      string
	FS        *vaults.FileSystem
	OnOpenCAD func(path string) // called when a “file” (incl. numeric-dir) is opened

	// UI
	tree       *widget.Tree
	nameLabel  *widget.Label
	pathEntry  *widget.Entry
	sizeLabel  *widget.Label
	timeLabel  *widget.Label
	refreshBtn *widget.Button

	// UI actions
	renameBtn *widget.Button
	moveBtn   *widget.Button
	copyBtn   *widget.Button
	delBtn    *widget.Button

	// state
	lastClickPath string
	lastClickAt   time.Time
	selectedUID   string

	// inside type VaultTab
	infoCache map[string]vaults.FileInfo // abs-node-id -> FileInfo
}

func NewVaultTab(win fyne.Window, root string, onOpenCAD func(string)) *VaultTab {
	vt := &VaultTab{
		Root:      root,
		OnOpenCAD: onOpenCAD,
		infoCache: make(map[string]vaults.FileInfo),
	}

	userName, err := user.Current()
	if err != nil {
		log.Fatalf("user.Current() error: %s", err)
	}

	vt.FS, err = vaults.NewClientFileSystem(root, userName.Username)
	if err != nil {
		log.Fatalf("initialization failed, %v", err)
	}

	vt.tree = widget.NewTree(
		// children
		func(uid widget.TreeNodeID) []widget.TreeNodeID {
			// Log every children request so you can see activity
			log.Printf("children(uid=%q)", string(uid))
			if uid == "" {
				// Expose exactly one root node
				return []widget.TreeNodeID{widget.TreeNodeID(vt.Root)}
			}
			return toIDs(vt.listChildren(string(uid))) // <-- use vt.listChildren (method)
		},
		// isBranch? (directory?)
		func(uid widget.TreeNodeID) bool {
			if uid == "" { // virtual root
				return true
			}
			// Root should behave as a branch so it can open
			if string(uid) == vt.Root {
				return true
			}
			// Ask the parent listing and find this entry
			if fi, ok := vt.lookupInfo(string(uid)); ok {
				return fi.IsDir()
			}
			return false
		},
		// create node UI (icon + label)
		func(branch bool) fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FolderIcon()),
				widget.NewLabel("..."),
			)
		},
		// update node UI
		func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			row := obj.(*fyne.Container)
			icon := row.Objects[0].(*widget.Icon)
			lbl := row.Objects[1].(*widget.Label)

			if uid == "" {
				lbl.SetText(filepath.Base(vt.Root))
				icon.SetResource(theme.FolderOpenIcon())
				return
			}

			base := filepath.Base(string(uid))
			lbl.SetText(base)

			if string(uid) == vt.Root {
				// Show opened-folder icon for the root
				icon.SetResource(theme.FolderOpenIcon())
				return
			}
			if fi, ok := vt.lookupInfo(string(uid)); ok {
				if fi.IsDir() {
					icon.SetResource(theme.FolderIcon())
				} else {
					icon.SetResource(theme.FileIcon())
				}
			} else {
				icon.SetResource(theme.WarningIcon())
			}
		},
	)

	vt.tree.OpenBranch(widget.TreeNodeID(vt.Root))
	vt.tree.Select(widget.TreeNodeID(vt.Root))
	vt.tree.Refresh()

	// Select behavior: single-click just updates details
	vt.tree.OnSelected = func(uid widget.TreeNodeID) {
		vt.selectedUID = string(uid)
		vt.showDetails(string(uid))
		vt.renameBtn.Enable()
		vt.moveBtn.Enable()
		vt.copyBtn.Enable()
		vt.delBtn.Enable()

		now := time.Now()
		if vt.lastClickPath == string(uid) && now.Sub(vt.lastClickAt) <= 300*time.Millisecond {
			vt.openPath(win, string(uid))
			vt.tree.Unselect(uid)
			vt.selectedUID = ""
			vt.lastClickPath, vt.lastClickAt = "", time.Time{}
			return
		}
		vt.lastClickPath, vt.lastClickAt = string(uid), now
	}

	vt.tree.OnUnselected = func(uid widget.TreeNodeID) {
		if vt.selectedUID == string(uid) {
			vt.selectedUID = ""
			vt.renameBtn.Disable()
			vt.moveBtn.Disable()
			vt.copyBtn.Disable()
			vt.delBtn.Disable()
		}
	}

	// Key: Enter opens selected
	if c, ok := win.Canvas().(interface {
		AddShortcut(shortcut fyne.Shortcut, handler func(fyne.Shortcut))
	}); ok {
		c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyReturn}, func(fyne.Shortcut) {
			uid := vt.currentSelection()
			if uid == "" {
				return
			}
			vt.openPath(win, uid)
			vt.tree.Unselect(uid)
			vt.selectedUID = ""
		})
		c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyEnter}, func(fyne.Shortcut) {
			uid := vt.currentSelection()
			if uid == "" {
				return
			}
			vt.openPath(win, uid)
			vt.tree.Unselect(uid)
		})
	}

	vt.renameBtn = widget.NewButtonWithIcon("Rename", theme.DocumentCreateIcon(), func() {
		uid := vt.selectedUID
		if uid == "" {
			return
		}
		entry := widget.NewEntry()
		entry.SetText(filepath.Base(uid))
		dialog.ShowForm("Rename", "OK", "Cancel", []*widget.FormItem{
			widget.NewFormItem("New name", entry),
		}, func(ok bool) {
			if !ok {
				return
			}
			newName := strings.TrimSpace(entry.Text)
			if newName == "" || newName == filepath.Base(uid) {
				return
			}
			if err := renamePath(uid, filepath.Join(filepath.Dir(uid), newName)); err != nil {
				dialog.ShowError(err, win)
				return
			}
			vt.refreshTree()
		}, win)
	})

	vt.moveBtn = widget.NewButtonWithIcon("Move", theme.NavigateNextIcon(), func() {
		uid := vt.selectedUID
		if uid == "" {
			return
		}
		dest := widget.NewEntry()
		dest.SetPlaceHolder("Destination directory")
		dialog.ShowForm("Move", "OK", "Cancel", []*widget.FormItem{
			widget.NewFormItem("To dir", dest),
		}, func(ok bool) {
			if !ok {
				return
			}
			d := strings.TrimSpace(dest.Text)
			if d == "" {
				return
			}
			if err := movePath(uid, filepath.Join(d, filepath.Base(uid))); err != nil {
				dialog.ShowError(err, win)
				return
			}
			vt.refreshTree()
		}, win)
	})

	vt.copyBtn = widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		uid := vt.selectedUID
		if uid == "" {
			return
		}
		dest := widget.NewEntry()
		dest.SetPlaceHolder("Destination directory")
		dialog.ShowForm("Copy", "OK", "Cancel", []*widget.FormItem{
			widget.NewFormItem("To dir", dest),
		}, func(ok bool) {
			if !ok {
				return
			}
			d := strings.TrimSpace(dest.Text)
			if d == "" {
				return
			}
			if err := copyPath(uid, filepath.Join(d, filepath.Base(uid))); err != nil {
				dialog.ShowError(err, win)
				return
			}
			vt.refreshTree()
		}, win)
	})

	vt.delBtn = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		uid := vt.selectedUID
		if uid == "" {
			return
		}
		dialog.ShowConfirm("Delete", "Are you sure?\n"+uid, func(yes bool) {
			if !yes {
				return
			}
			if err := os.RemoveAll(uid); err != nil {
				dialog.ShowError(err, win)
				return
			}
			vt.selectedUID = ""
			vt.refreshTree()
		}, win)
	})

	c := win.Canvas()
	c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF2}, func(fyne.Shortcut) {
		if vt.selectedUID != "" {
			vt.renameBtn.Tapped(nil)
		}
	})
	c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyDelete}, func(fyne.Shortcut) {
		if vt.selectedUID != "" {
			vt.delBtn.Tapped(nil)
		}
	})
	c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF5}, func(fyne.Shortcut) {
		vt.refreshTree()
	})

	// standard off until selection
	vt.renameBtn.Disable()
	vt.moveBtn.Disable()
	vt.copyBtn.Disable()
	vt.delBtn.Disable()

	// Right: details + toolbar
	vt.nameLabel = widget.NewLabel("")
	vt.nameLabel.Wrapping = fyne.TextWrapWord
	vt.pathEntry = widget.NewEntry()
	vt.pathEntry.Disable()
	vt.sizeLabel = widget.NewLabel("")
	vt.timeLabel = widget.NewLabel("")
	vt.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		vt.refreshTree()
	})

	details := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Name:"), vt.nameLabel,
		widget.NewLabel("Path:"), vt.pathEntry,
		widget.NewLabel("Size:"), vt.sizeLabel,
		widget.NewLabel("Modified:"), vt.timeLabel,
		widget.NewSeparator(),
		container.NewHBox(vt.refreshBtn, vt.renameBtn, vt.moveBtn, vt.copyBtn, vt.delBtn),
	)

	split := container.NewHSplit(vt.tree, details)
	split.Offset = 0.45

	vt.Tab = container.NewTabItemWithIcon(filepath.Base(root), theme.FolderOpenIcon(), split)
	return vt
}

// ---------- Tree helpers ----------

// listChildren returns immediate children as absolute node IDs.
// It converts the absolute 'abs' into a vault-relative path for vfs.ListDir.
func (vt *VaultTab) listChildren(abs string) []string {
	rel := vt.rel(abs)                 // "" for root
	entries, err := vt.FS.ListDir(rel) // your FS expects relative paths
	if err != nil {
		log.Printf("ListDir(%q) error: %v", rel, err)
		return nil
	}
	log.Printf("ListDir(%q) -> %d entries", rel, len(entries))

	out := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		// If your FS already filters dot-items, you can drop this:
		if strings.HasPrefix(name, ".") {
			continue
		}
		out = append(out, filepath.Join(abs, name)) // absolute node IDs
	}
	return out
}

// rel returns the vault-relative path for an absolute path under vt.Root.
// "" means the vault root itself.
func (vt *VaultTab) rel(abs string) string {
	if abs == "" {
		return ""
	}
	rootClean, _ := filepath.Abs(vt.Root)
	absClean, _ := filepath.Abs(abs)

	if absClean == rootClean {
		return ""
	}
	pref := rootClean + string(os.PathSeparator)
	if strings.HasPrefix(absClean, pref) {
		return strings.TrimPrefix(absClean, pref)
	}
	return "" // outside → treat as root
}

func (vt *VaultTab) currentSelection() string {
	return vt.selectedUID
}

func (vt *VaultTab) refreshTree() {
	// Force refresh of the currently visible nodes
	vt.tree.Refresh()
}

// toIDs converts []string node IDs into []widget.TreeNodeID
func toIDs(paths []string) []widget.TreeNodeID {
	out := make([]widget.TreeNodeID, len(paths))
	for i, p := range paths {
		out[i] = widget.TreeNodeID(p)
	}
	return out
}

// ---------- Details & open behavior ----------

func (vt *VaultTab) showDetails(path string) {
	vt.nameLabel.SetText(filepath.Base(path))
	vt.pathEntry.SetText(path)
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() {
			vt.sizeLabel.SetText("(directory)")
		} else {
			vt.sizeLabel.SetText(fmt.Sprintf("%d bytes", fi.Size()))
		}
		vt.timeLabel.SetText(fi.ModTime().Format(time.RFC3339))
	} else {
		vt.sizeLabel.SetText("unknown")
		vt.timeLabel.SetText("-")
	}
}

func (vt *VaultTab) openPath(win fyne.Window, path string) {
	fi, err := os.Stat(path)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	if fi.IsDir() {
		name := filepath.Base(path)
		if reNumeric.MatchString(name) {
			// "numeric dir" behaves like a file-container
			if vt.OnOpenCAD != nil {
				vt.OnOpenCAD(path)
				return
			}
			dialog.ShowInformation("Open", "Open numeric file container:\n"+path, win)
			return
		}
		// ordinary folder: toggle open/close in the tree
		id := widget.TreeNodeID(path)
		if vt.tree.IsBranchOpen(id) {
			vt.tree.CloseBranch(id)
		} else {
			vt.tree.OpenBranch(id)
		}
		return
	}

	// real file
	if vt.OnOpenCAD != nil {
		vt.OnOpenCAD(path)
		return
	}
	dialog.ShowInformation("Open", "Open file:\n"+path, win)
}

// ---------- FS ops (safe-ish) ----------

func renamePath(src, dst string) error {
	if src == "" || dst == "" {
		return errors.New("empty path")
	}
	if sameFile(src, dst) {
		return nil
	}
	return os.Rename(src, dst)
}

func movePath(src, dst string) error {
	// Try rename (cheap); if cross-device, fallback to copy+delete.
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if !isCrossDevice(err) {
		return err
	}
	if err := copyPath(src, dst); err != nil {
		return err
	}
	return os.RemoveAll(src)
}

func copyPath(src, dst string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst, fi.Mode())
}

func copyFile(src, dst string, mode fs.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	_, cErr := io.Copy(out, in)
	cErr2 := out.Close()
	if cErr != nil {
		return cErr
	}
	return cErr2
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return copyFile(path, target, info.Mode())
	})
}

func sameFile(a, b string) bool {
	aa, _ := filepath.Abs(a)
	bb, _ := filepath.Abs(b)
	return aa == bb
}

func isCrossDevice(err error) bool {
	// crude check: on unix, EXDEV; windows: use message
	if errors.Is(err, os.ErrInvalid) {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "invalid cross-device link") || strings.Contains(strings.ToLower(msg), "exdev")
}

// lookupInfo returns FileInfo for a given absolute node id.
// It first checks the cache, then falls back to listing the parent once.
func (vt *VaultTab) lookupInfo(abs string) (vaults.FileInfo, bool) {
	if abs == "" || abs == "." || abs == vt.Root {
		return vaults.FileInfo{}, false
	}

	if fi, ok := vt.infoCache[abs]; ok {
		return fi, true
	}

	parentAbs := filepath.Dir(abs)
	relParent := vt.rel(parentAbs)
	entries, err := vt.FS.ListDir(relParent)
	if err != nil {
		log.Printf("lookupInfo: ListDir(%q) error: %v", relParent, err)
		return vaults.FileInfo{}, false
	}
	for _, fi := range entries {
		childAbs := filepath.Join(parentAbs, fi.Name())
		vt.infoCache[childAbs] = fi
		if childAbs == abs {
			return fi, true
		}
	}
	log.Printf("lookupInfo: %q not found in parent %q", abs, parentAbs)
	return vaults.FileInfo{}, false
}
