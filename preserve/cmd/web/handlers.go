package main

import (
	"fmt"
	"net/http"
	"preserve/pkg/models"
	"strconv"
	"text/template"
)

// Function to handle main homepage.
func (preserve *Preserve) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		preserve.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Parsing the template files.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		preserve.logger.Error("Faced an Error in Parsing Template: %s", err.Error())
		preserve.serverError(w, err)
	}
	// Executing the templates.
	err = ts.Execute(w, nil)
	if err != nil {
		preserve.logger.Error("Faced an Error while writing Response: %s", err.Error())
		preserve.serverError(w, err)
	}
}

// Function to show a note.
func (preserve *Preserve) showNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		preserve.notFound(w)
		return
	}

	note, err := preserve.notes.Get(id)
	if err == models.ErrNoRecord {
		preserve.notFound(w)
		return
	} else if err != nil {
		preserve.serverError(w, err)
		return
	}

	data := &templateData{Note: note}
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		preserve.serverError(w, err)
	}
	err = ts.Execute(w, data)
	if err != nil {
		preserve.serverError(w, err)
	}
}

// Function to create a note.
func (preserve *Preserve) createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		preserve.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "0 snail"
	content := "0 snail\n Climb Mount Fuji,\nBut slowly, slowly."
	expires := "7"
	id, err := preserve.notes.Insert(title, content, expires)
	if err != nil {
		preserve.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/note?id=%d", id), http.StatusSeeOther)
}
