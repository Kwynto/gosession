package main

import (
	"html/template"
	"path/filepath"
)

// Structure for the data template
type templateData struct {
	User        string
	Hash        string
	Cart        []string
	Transitions []string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// We use the filepath.Glob function to get a slice of all file paths with the extension '.page.tmpl'.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// We iterate through the template file from each page.
	for _, page := range pages {
		// Extracting the final file name
		name := filepath.Base(page)

		// Processing the iterated template file.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// We use the ParseGlob method to add all the wireframe templates.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// We use the ParseGlob method to add all auxiliary templates.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Adding the resulting set of templates to the cache using the page name
		cache[name] = ts
	}

	// We return the received map.
	return cache, nil
}
