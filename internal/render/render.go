package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/calvarado2004/bookings/internal/config"
	"github.com/calvarado2004/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

var pathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	app = a
}

//AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()

	}

	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("could not get template from template cache")

		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil

}

//Map of functions available to the templates
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

//HumanDate returns time in yyyy-mm-dd format
func HumanDate(t time.Time) string {

	return t.Format("2006-01-02")
}

//FormatDate helps to format numeric date into string date
func FormatDate(t time.Time, f string) string {

	return t.Format(f)
}

//Add add two numbers
func Add(a, b int) int {

	return a + b
}

//Iterate performs a for loop
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}

	return items
}

// Creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}
