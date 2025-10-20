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
	"github.com/grd/FreePDM/internal/vault/localfs"
)

var reNumeric = regexp.MustCompile(`^\d+$`)

type VaultTab struct {
	Tab       *container.TabItem
	Root      string
	FS        *localfs.FileSystem
	OnOpenCAD func(path string) // called when a “file” (incl. numeric-dir) is opened
	win       fyne.Window

	// UI
	tree       *widget.Tree
	nameLabel  *widget.Label
	pathEntry  *widget.Entry
	sizeLabel  *widget.Label
	timeLabel  *widget.Label
	refreshBtn *widget.Button

	// UI actions
	renameBtn   *widget.Button
	moveBtn     *widget.Button
	copyBtn     *widget.Button
	delBtn      *widget.Button
	allocateBtn *widget.Button
	assignBtn   *widget.Button

	// state
	lastClickPath string
	lastClickAt   time.Time
	selectedUID   string

	// inside type VaultTab
	infoCache map[string]localfs.FileInfo // abs-node-id -> FileInfo
}

// NewVaultTab builds the Vault UI tab.
//   - Shows a tree on the left and a details/toolbar pane on the right.
//   - "Allocate" is visible only when a regular (non-numeric) directory is selected.
//   - "Assign" is visible only when a numeric container is selected that still
//     contains the placeholder file at version 0 (i.e.,  "<container>/0/<EmptyFile>").
//
// Note: Only containers are read-only in your policy; normal directories are writable.
// Allocate relies on fs.Allocate(dstDir). Assign relies on fs.Assign(containerNumber, fileName).
func NewVaultTab(win fyne.Window, root string, onOpenCAD func(string)) *VaultTab {
	vt := &VaultTab{
		Root:      root,
		OnOpenCAD: onOpenCAD,
		infoCache: make(map[string]localfs.FileInfo),
		win:       win,
	}

	fs, err := localfs.NewClientFileSystem(root)
	if err != nil {
		log.Fatalf("initialization failed, %v", err)
	}
	vt.FS = fs

	// --- Tree setup -----------------------------------------------------------
	vt.tree = widget.NewTree(
		// children
		func(uid widget.TreeNodeID) []widget.TreeNodeID {
			if uid == "" {
				// Expose exactly one root node
				return []widget.TreeNodeID{widget.TreeNodeID(vt.Root)}
			}
			return toIDs(vt.listChildren(string(uid)))
		},
		// isBranch? (directory?)
		func(uid widget.TreeNodeID) bool {
			if uid == "" { // virtual root
				return true
			}
			if string(uid) == vt.Root {
				return true
			}
			if fi, ok := vt.lookupInfo(string(uid)); ok {
				return fi.IsDir()
			}
			return false
		},

		// create node UI (main icon + label + right-aligned status badge)
		func(branch bool) fyne.CanvasObject {
			icon := widget.NewIcon(theme.FolderIcon())
			lbl := widget.NewLabel("")

			badge := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {})
			badge.Importance = widget.LowImportance
			badge.Hide()

			left := container.NewHBox(icon, lbl)
			row := container.NewBorder(nil, nil, left, badge, nil)
			return row
		},

		func(uid widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			row := obj.(*fyne.Container)             // Border
			left := row.Objects[0].(*fyne.Container) // HBox(icon,label)
			badge := row.Objects[1].(*widget.Button) // right-side badge

			icon := left.Objects[0].(*widget.Icon)
			lbl := left.Objects[1].(*widget.Label)

			abs := string(uid)

			if fi, ok := vt.lookupInfo(abs); ok {
				name := fi.Name()
				if fi.Alloc() == localfs.AllocAllocatedWithCandidate {
					if cand, ok := fi.AllocCandidate(); ok && cand != "" {
						name = cand
					}
				}
				if name == "" {
					name = "(unnamed)" // guard
				}
				lbl.SetText(name)

				if fi.IsDir() {
					icon.SetResource(theme.FolderIcon())
				} else {
					icon.SetResource(theme.FileIcon())
				}

				switch fi.Alloc() {
				case localfs.AllocAllocatedEmpty:
					badge.SetIcon(theme.ConfirmIcon())
					badge.OnTapped = func() {
						dialog.ShowInformation(
							"Status",
							"Allocated — empty container (version 0). Use Assign to name it.",
							win,
						)
					}
					badge.Show()

				case localfs.AllocAllocatedWithCandidate:
					cand, _ := fi.AllocCandidate()
					badge.SetIcon(theme.DocumentCreateIcon())
					badge.OnTapped = func() {
						msg := "Unassigned — candidate found in version 0"
						if cand != "" {
							msg += ": “" + cand + "”"
						}
						msg += ". Use Assign to finalize."
						dialog.ShowInformation("Status", msg, win)
					}
					badge.Show()

				default:
					badge.Hide()
				}
				return
			}

			lbl.SetText(filepath.Base(abs))
			icon.SetResource(theme.FileIcon())
			badge.Hide()
		},
	)

	vt.tree.OpenBranch(widget.TreeNodeID(vt.Root))
	vt.tree.Select(widget.TreeNodeID(vt.Root))
	vt.tree.Refresh()

	// --- Buttons: Rename/Move/Copy/Delete ------------------------------------
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
			if err := renamePath(vt, uid, filepath.Join(filepath.Dir(uid), newName)); err != nil {
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

	// --- Buttons: Refresh / Allocate / Assign ---------------------------------
	vt.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
		vt.refreshTree()
	})

	// Allocate: creates an empty container in the selected directory (non-numeric dir)
	vt.allocateBtn = widget.NewButtonWithIcon("Allocate", theme.ContentAddIcon(), func() {
		abs := vt.selectedUID
		if abs == "" {
			dialog.ShowInformation("Allocate", "Select a directory first.", win)
			return
		}
		if fi, ok := vt.lookupInfo(abs); ok && !fi.IsDir() {
			abs = filepath.Dir(abs)
		}

		rel := vt.rel(abs)

		fl, err := vt.FS.Allocate(rel)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		vt.refreshTree()
		absContainer := filepath.Join(vt.Root, fl.Path, fl.ContainerNumber)
		vt.tree.Select(widget.TreeNodeID(absContainer))
	})

	// Assign: only relevant for an allocated container that still has the EmptyFile in version 0.
	vt.assignBtn = widget.NewButtonWithIcon("Assign", theme.DocumentCreateIcon(), func() {
		abs := vt.selectedUID
		if abs == "" {
			return
		}
		cn := filepath.Base(abs) // container number (numeric string)
		// Prompt for a human-friendly filename
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter filename.ext")
		dialog.ShowForm("Assign name", "Assign", "Cancel",
			[]*widget.FormItem{widget.NewFormItem("File name", entry)},
			func(ok bool) {
				if !ok {
					return
				}
				name := strings.TrimSpace(entry.Text)
				if name == "" {
					return
				}
				if err := vt.FS.Assign(cn, name); err != nil {
					dialog.ShowError(err, win)
					return
				}
				vt.refreshTree()
			},
			win,
		)
	})

	// Initially disable/hide action buttons until a selection is made.
	vt.renameBtn.Disable()
	vt.moveBtn.Disable()
	vt.copyBtn.Disable()
	vt.delBtn.Disable()

	vt.allocateBtn.Disable()
	vt.allocateBtn.Hide()

	vt.assignBtn.Disable()
	vt.assignBtn.Hide()

	// --- Details pane (right side) --------------------------------------------
	vt.nameLabel = widget.NewLabel("")
	vt.nameLabel.Wrapping = fyne.TextWrapWord
	vt.pathEntry = widget.NewEntry()
	vt.pathEntry.Disable()
	vt.sizeLabel = widget.NewLabel("")
	vt.timeLabel = widget.NewLabel("")

	details := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Name:"), vt.nameLabel,
		widget.NewLabel("Path:"), vt.pathEntry,
		widget.NewLabel("Size:"), vt.sizeLabel,
		widget.NewLabel("Modified:"), vt.timeLabel,
		widget.NewSeparator(),
		// Toolbar: Refresh + Allocate + Assign + (Rename/Move/Copy/Delete)
		container.NewHBox(vt.refreshBtn, vt.allocateBtn, vt.assignBtn, vt.renameBtn, vt.moveBtn, vt.copyBtn, vt.delBtn),
	)

	vt.tree.OnSelected = func(uid widget.TreeNodeID) {
		vt.selectedUID = string(uid)
		vt.showDetails(string(uid))

		// Algemene acties
		vt.renameBtn.Enable()
		vt.moveBtn.Enable()
		vt.copyBtn.Enable()
		vt.delBtn.Enable()

		path := string(uid)
		base := filepath.Base(path)

		// Allocate: only on regular (non-numeric) directories
		if fi, ok := vt.lookupInfo(path); ok && fi.IsDir() && !reNumeric.MatchString(base) {
			vt.allocateBtn.Show()
			vt.allocateBtn.Enable()
		} else {
			vt.allocateBtn.Hide()
			vt.allocateBtn.Disable()
		}

		// Assign: visible on containers with alloc-status (AllocatedEmpty / AllocatedWithCandidate)
		if fi, ok := vt.lookupInfo(path); ok &&
			!fi.IsDir() &&
			(fi.Alloc() == localfs.AllocAllocatedEmpty || fi.Alloc() == localfs.AllocAllocatedWithCandidate) {

			vt.assignBtn.Show()
			vt.assignBtn.Enable()

			vt.assignBtn.OnTapped = func() {
				// Resolve numeric container number from FileInfo (never from the UI label)
				cn := fi.ContainerNumber()
				if cn == "" {
					dialog.ShowError(fmt.Errorf("unable to resolve container number for %q", fi.Name()), win)
					return
				}

				// Prefill with candidate if available (only from <container>/0/)
				entry := widget.NewEntry()
				if cand, ok := fi.AllocCandidate(); ok && cand != "" {
					entry.SetText(cand)
				} else {
					entry.SetPlaceHolder("filename.ext")
				}

				dialog.ShowForm(
					"Assign name",
					"Assign",
					"Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("File name", entry),
					},
					func(ok bool) {
						if !ok {
							return
						}
						name := strings.TrimSpace(entry.Text)
						if name == "" {
							dialog.ShowError(fmt.Errorf("file name is required"), win)
							return
						}
						if err := vt.FS.Assign(cn, name); err != nil {
							dialog.ShowError(err, win)
							return
						}
						vt.refreshTree()
					},
					win,
				)
			}
		} else {
			vt.assignBtn.Hide()
			vt.assignBtn.Disable()
		}

		// Double-click to open branch
		now := time.Now()
		if vt.lastClickPath == path && now.Sub(vt.lastClickAt) <= 300*time.Millisecond {
			vt.openPath(win, path)
			vt.tree.Unselect(uid)
			vt.selectedUID = ""
			vt.lastClickPath, vt.lastClickAt = "", time.Time{}
			return
		}
		vt.lastClickPath, vt.lastClickAt = path, now
	}

	vt.tree.OnUnselected = func(uid widget.TreeNodeID) {
		if vt.selectedUID == string(uid) {
			vt.selectedUID = ""
			vt.renameBtn.Disable()
			vt.moveBtn.Disable()
			vt.copyBtn.Disable()
			vt.delBtn.Disable()

			vt.allocateBtn.Hide()
			vt.allocateBtn.Disable()

			vt.assignBtn.Hide()
			vt.assignBtn.Disable()
		}
	}

	// --- Keyboard shortcuts ----------------------------------------------------
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
	}

	// --- Layout ---------------------------------------------------------------
	split := container.NewHSplit(vt.tree, details)
	split.Offset = 0.45
	vt.Tab = container.NewTabItemWithIcon(filepath.Base(root), theme.FolderOpenIcon(), split)
	return vt
}

// ---------- Tree helpers ----------

// listChildren returns immediate children as absolute node IDs.
// It converts the absolute 'abs' into a vault-relative path for localfs.ListDir.
func (vt *VaultTab) listChildren(abs string) []string {
	rel := vt.rel(abs)                 // "" for root
	entries, err := vt.FS.ListDir(rel) // your FS expects relative paths
	if err != nil {
		log.Printf("ListDir(%q) error: %v", rel, err)
		return nil
	}

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

func renamePath(vt *VaultTab, src, dst string) error {
	if src == "" || dst == "" {
		return errors.New("empty path")
	}
	if sameFile(src, dst) {
		return nil
	}
	ok, err := localfs.IsEmptyDirectory(src)
	if err != nil {
		log.Print("hello3")
		return err
	}
	if ok {
		return os.Rename(src, dst)
	}
	ok, err = localfs.IsContainer(src)
	if err != nil {
		log.Print("hello2")
		return err
	}
	if ok {
		log.Print("hello")
		return vt.FS.FileRename(src, dst)
	}

	// No empty, no Container, it must be a directory rename
	return vt.FS.DirectoryRename(src, dst)
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
func (vt *VaultTab) lookupInfo(abs string) (localfs.FileInfo, bool) {
	if abs == "" || abs == "." || abs == vt.Root {
		return localfs.FileInfo{}, false
	}

	if fi, ok := vt.infoCache[abs]; ok {
		return fi, true
	}

	parentAbs := filepath.Dir(abs)
	relParent := vt.rel(parentAbs)
	entries, err := vt.FS.ListDir(relParent)
	if err != nil {
		return localfs.FileInfo{}, false
	}
	for _, fi := range entries {
		childAbs := filepath.Join(parentAbs, fi.Name())
		vt.infoCache[childAbs] = fi
		if childAbs == abs {
			return fi, true
		}
	}
	return localfs.FileInfo{}, false
}

func (vt *VaultTab) Window() fyne.Window { return vt.win }
