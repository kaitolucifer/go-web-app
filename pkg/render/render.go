package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/kaitolucifer/go-web-app/pkg/config"
	"github.com/kaitolucifer/go-web-app/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData add default data to templates
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTamplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	// get the template cache from the app config
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, AddDefaultData(td))

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

//CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return tmplCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return tmplCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = ts
	}

	return tmplCache, nil
}
