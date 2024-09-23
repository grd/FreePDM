// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	// "fmt"
	"html/template"
	// "log"
	"net/http"
	"os"
	"path"
	// "github.com/grd/FreePDM/src/config"
	// ex "github.com/grd/FreePDM/src/extras"
)

// Get path to template directory from env variable
var templatesDir = os.Getenv("TEMPLATES_DIR")

// loveMountains renders the love-mountains page after passing some data to the HTML template
func loveMountains(w http.ResponseWriter, r *http.Request) {
	// Build path to template
	tmplPath := path.Join(templatesDir, "love-mountains.html")
	// Load template from disk
	tmpl := template.Must(template.ParseFiles(tmplPath))
	// Inject data into template
	data := "La Chartreuse"
	tmpl.Execute(w, data)
}

func main() {
	// Create route to love-mountains web page
	http.HandleFunc("/love-mountains", loveMountains)
	// Launch web server on port 80
	http.ListenAndServe(":80", nil)
}
