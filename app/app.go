package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbname string) {
	var err error
	database, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS post (id STRING PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/posts", a.getPosts).Methods("GET")
	a.Router.HandleFunc("/posts/{id}", a.getPost).Methods("GET")
	a.Router.HandleFunc("/posts", a.createPost).Methods("POST")
	a.Router.HandleFunc("/posts/{id}", a.deletePost).Methods("DELETE")
}

func (a *App) getPosts(w http.ResponseWriter, r *http.Request) {

}
func (a *App) getPost(w http.ResponseWriter, r *http.Request) {

}
func (a *App) createPost(w http.ResponseWriter, r *http.Request) {

}
func (a *App) deletePost(w http.ResponseWriter, r *http.Request) {

}
