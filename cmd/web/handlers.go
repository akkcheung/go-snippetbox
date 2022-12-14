package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/akkcheung/go-snippetbox/pkg/forms"
	"github.com/akkcheung/go-snippetbox/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//debug
	//fmt.Printf("-> home")

	app.render(w, r, "home.page.tmpl", nil)

}

func (app *application) all(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}

	app.render(w, r, "all.page.tmpl", data)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//debug
	//fmt.Printf("%s", vars["id"])

	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
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

	//flash := app.session.PopString(r, "flash")

	//debug
	//fmt.Printf("flash -> %s", flash)

	app.render(w, r, "show.page.tmpl", &templateData{
		//Flash:   flash,
		Snippet: s,
	})

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	//debug
	fmt.Printf("-> createSnippet")

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet successfully created")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	//app.render(w, r, "create.page.tmpl", nil)
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
