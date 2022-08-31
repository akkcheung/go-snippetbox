package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	s := r.PathPrefix("/snippet").Subrouter()
	s.Use(secureHeaders, app.logRequest, app.recoverPanic)

	r.HandleFunc("/", app.home).Methods("GET")

	s.HandleFunc("/all", app.all).Methods("GET")
	s.HandleFunc("/create", app.createSnippetForm).Methods("GET")
	s.HandleFunc("/create", app.createSnippet).Methods("POST")
	s.HandleFunc("/{id}", app.showSnippet).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))

	r.Use(secureHeaders, app.logRequest, app.recoverPanic, app.session.Enable)
	s.Use(app.session.Enable)

	return r
}
