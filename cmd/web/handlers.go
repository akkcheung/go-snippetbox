package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akkcheung/go-snippetbox/pkg/models"
)

//func home(w http.ResponseWriter, r *http.Request) {
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		app.notFound(w)
		return
	}

	//w.Write([]byte("Hello from Snippetbox"))

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v/n", snippet)
	}

	/*
		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}
	*/

	//ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	/*
		ts, err := template.ParseFiles(files...)
		if err != nil {
	*/

	//log.Println(err.Error())

	//app.errorLog.Println(err.Error())
	//http.Error(w, "Internal Server Error", 500)

	/*
			app.serverError(w, err)
			return
		}
	*/

	/*
		err = ts.Execute(w, nil)
		if err != nil {
	*/
	//log.Println(err.Error())

	//app.errorLog.Println(err.Error())
	//http.Error(w, "Internal Server Error", 500)

	/*
			app.serverError(w, err)
		}
	*/
}

//func showSnippet(w http.ResponseWriter, r *http.Request) {
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		//http.NotFound(w, r)
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	fmt.Fprintf(w, "%v", s)
}

//func createSnippet(w http.ResponseWriter, r *http.Request) {
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		//		w.WriteHeader(405)
		//		w.Write([]byte("Method Not Allowed"))
		// http.Error(w, "Method Not Allowed", 405)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	//w.Write([]byte("Create a new snippet..."))

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
