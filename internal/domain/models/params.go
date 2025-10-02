// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

// Parameters holds structured metadata for a specific file version.
type Parameters struct {
	VersionID   string         // linked version
	Data        map[string]any // schema-less key→value(s)
	IndexedKeys []string       // keys commonly filtered on
}

// Range for numeric/date filters in search queries.
type Range struct {
	Min any
	Max any
}

// SearchQuery expresses faceted/full-text queries over parameters.
type SearchQuery struct {
	Vault   string
	Text    string
	Filters map[string][]string // exact filters: key → multiple values
	Ranges  map[string]Range    // per-key range constraints
	Limit   int
	Offset  int
}

// SearchResult returns a navigable hit in a vault.
type SearchResult struct {
	Vault      string
	RelPath    string
	VersionID  string
	Score      float64           // optional relevance
	Highlights map[string]string // key → snippet
}
