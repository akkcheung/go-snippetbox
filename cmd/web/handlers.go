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

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}

	app.render(w, r, "home.page.tmpl", data)
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

	app.render(w, r, "show.page.tmpl", &templateData{
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

	//	title := r.PostForm.Get("title")
	//	content := r.PostForm.Get("content")
	//	expires := r.PostForm.Get("expires")
	//
	//	errors := make(map[string]string)
	//
	//	if strings.TrimSpace(title) == "" {
	//		errors["title"] = "This field cannot be blank"
	//	} else if utf8.RuneCountInString(title) > 100 {
	//		errors["title"] = "This field is too long (maximum is 100 characters"
	//	}
	//
	//	if strings.TrimSpace(content) == "" {
	//		errors["content"] = "This field cannot be blank"
	//	}
	//
	//	if strings.TrimSpace(expires) == "" {
	//		errors["expires"] = "This field cannot be blank"
	//	} else if expires != "365" && expires != "7" && expires != "1" {
	//		errors["expires"] = "This field is invalid"
	//	}
	//
	//	if len(errors) > 0 {
	//		app.render(w, r, "create.page.tmpl", &templateData{
	//			FormErrors: errors,
	//			FormData:   r.PostForm,
	//		})
	//		return
	//	}
	//

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	//id, err := app.snippets.Insert(title, content, expires)
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	//app.render(w, r, "create.page.tmpl", nil)
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
