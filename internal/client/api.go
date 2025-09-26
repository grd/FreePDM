// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package client

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type API struct {
	BaseURL string
	HTTP    *http.Client
}

func (a *API) ListVaults() (any, any) {
	panic("unimplemented")
}

func New(base string) *API {
	jar, _ := cookiejar.New(nil)
	return &API{
		BaseURL: base,
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
			Jar:     jar,
		},
	}
}

func (a *API) Login(user, pass string) error {
	// TODO: implement POST /login and error handling.
	return nil
}

// Domain types
type Vault struct {
	Name string
	// add fields later (ID, path, etc.)
}
