package main

import (
	"html/template" 
    "path/filepath"
	"forum.christen.net/internal/models"
	"net/http"
	"time"
)

type templateData struct {
	CurrentYear int
	User models.User
	Thread models.Thread
	Threads []models.Thread
	Msg models.Message
	Form any
	Flash string
	IsAuthenticated bool
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}

    pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        name := filepath.Base(page)

		temp, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		temp, err = temp.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

        temp, err = temp.ParseFiles(page)
        if err != nil {
            return nil, err
        }

        cache[name] = temp
    }

    return cache, nil
}