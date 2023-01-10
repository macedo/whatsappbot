package handlers

import (
	"html/template"
	"io"
)

func renderPage(w io.Writer, data any) error {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
