package renders

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// functions is a map of template functions
var functions = template.FuncMap{}

var temps = make(map[string]*template.Template)

// getTemplateCache is a helper function to cache all HTML templates as a map
func getTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	baseDir := GetProjectRoot("views", "templates")

	templatesDir := filepath.Join(baseDir, "*.page.html")
	pages, err := filepath.Glob(templatesDir)
	if err != nil {

		return myCache, fmt.Errorf("error globbing templates: %v", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, fmt.Errorf("error parsing page %s: %v", name, err)
		}

		layoutsPath := filepath.Join(baseDir, "*.layout.html")
		matches, err := filepath.Glob(layoutsPath)
		if err != nil {
			return myCache, fmt.Errorf("error finding layout files: %v", err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(layoutsPath)
			if err != nil {
				return myCache, fmt.Errorf("error parsing layout files: %v", err)
			}
		}

		myCache[name] = ts

	}
	temps = myCache
	return myCache, nil
}

// GetProjectRoot dynamically finds the project root directory
func GetProjectRoot(first, second string) string {
	cwd, _ := os.Getwd()
	baseDir := cwd
	if strings.HasSuffix(baseDir, "cmd") {
		baseDir = filepath.Join(cwd, "../")
	}
	return filepath.Join(baseDir, first, second)
}
