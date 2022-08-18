package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/akkcheung/go-snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  errorLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	//mux := http.NewServeMux()

	//mux.HandleFunc("/", home)
	//mux.HandleFunc("/", app.home)

	//mux.HandleFunc("/snippet", showSnippet)
	//mux.HandleFunc("/snippet", app.showSnippet)

	//mux.HandleFunc("/snippet/create", createSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)

	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		//Handler:  mux,
		Handler: app.routes(),
	}

	//log.Println("Starting server on :4000")
	//log.Printf("Starting server on %s", *addr)
	infoLog.Printf("Starting server on %s", *addr)

	//err := http.ListenAndServe(":4000", mux)
	//err := http.ListenAndServe(*addr, mux)
	err = srv.ListenAndServe()

	//handler := cors.Default().Handler(mux)
	//err := http.ListenAndServe(":4000", handler)
	//log.Fatal(err)
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
