// Copyright 2025 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}
