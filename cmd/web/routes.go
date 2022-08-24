package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", app.home).Methods("GET")
	r.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
	r.HandleFunc("/snippet/{id}", app.showSnippet).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))

	return app.recoverPanic(app.logRequest(secureHeaders(r)))
}
