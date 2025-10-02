// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package params

import "github.com/grd/FreePDM/internal/domain/models"

// Store provides parameter access and search.
type Store interface {
	// Fetch parameters for latest or a specific version (if VersionID is non-empty).
	GetByPath(vault, rel, versionID string) (models.Parameters, error)

	// Faceted/full-text search over parameter space.
	Search(q models.SearchQuery) ([]models.SearchResult, error)
}
