package main

import (
	"html/template"
	"path/filepath"

	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"snippetbox.alexedwards.net/internal/models"
)

// Define a templateData type to act as the holding structure for any dynamic data that pass to HTML templates.
type templateData struct {
	Snippet *models.Snippet
	// Include ta Snippets field in the templateData struct.
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glop() function to get a slice of all filepaths that match the pattern "ui/html/pages/*.html".
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the file name from the path
		name := filepath.Base(page)

		// Parse the base template file
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name as key
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
