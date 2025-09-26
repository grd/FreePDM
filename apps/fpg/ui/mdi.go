// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type TabManager struct {
	Tabs *container.DocTabs
}

func NewTabManager() *TabManager {
	dt := container.NewDocTabs()
	dt.SetTabLocation(container.TabLocationTop)

	// Optional: veto closing certain tabs (e.g., "Home")
	// No dt.CloseIntercept

	// Optional: cleanup after a tab was closed
	dt.OnClosed = func(ti *container.TabItem) {
		// no-op; add cleanup if needed
	}

	return &TabManager{Tabs: dt}
}

// Add and select
func (m *TabManager) AddTab(title string, content fyne.CanvasObject) {
	it := container.NewTabItem(title, content)
	m.Tabs.Append(it)
	m.Tabs.Select(it)
}

// Close selected tab (uses DocTabs API)
func (m *TabManager) CloseSelected() {
	if it := m.Tabs.Selected(); it != nil {
		m.Tabs.Remove(it)
	}
}

// Navigation helpers
func (m *TabManager) MaybeSelectedIndex() int {
	cur := m.Tabs.Selected()
	for i, it := range m.Tabs.Items {
		if it == cur {
			return i
		}
	}
	return 0
}

func (m *TabManager) SelectNext() {
	n := len(m.Tabs.Items)
	if n == 0 {
		return
	}
	i := m.MaybeSelectedIndex()
	m.Tabs.Select(m.Tabs.Items[(i+1)%n])
}

func (m *TabManager) SelectPrev() {
	n := len(m.Tabs.Items)
	if n == 0 {
		return
	}
	i := m.MaybeSelectedIndex()
	m.Tabs.Select(m.Tabs.Items[(i-1+n)%n])
}
