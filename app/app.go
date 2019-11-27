package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/TheStevbeef/communityServiceGo/models"
	"github.com/TheStevbeef/communityServiceGo/utils"
	"github.com/rs/xid"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbname string) {
	var err error
	a.DB, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	statement, err := a.DB.Prepare("CREATE TABLE IF NOT EXISTS post (post_id TEXT PRIMARY KEY, timestamp TEXT, message TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	statement, err = a.DB.Prepare("CREATE TABLE IF NOT EXISTS user (user_id TEXT PRIMARY KEY, name TEXT, image_url TEXT, post_id TEXT, FOREIGN KEY(post_id) REFERENCES post(post_id))")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	statement, err = a.DB.Prepare("CREATE TABLE IF NOT EXISTS media (content_type TEXT, url TEXT, post_id TEXT, FOREIGN KEY(post_id) REFERENCES post(post_id))")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/posts", a.GetPosts).Methods("GET")
	a.Router.HandleFunc("/posts/{id}", a.GetPost).Methods("GET")
	a.Router.HandleFunc("/posts", a.CreatePost).Methods("POST")
	a.Router.HandleFunc("/posts/{id}", a.DeletePost).Methods("DELETE")
}

func (a *App) GetPosts(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Offset type: %T; Limit type: %T", offset, limit)
	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	products, err := models.GetPosts(a.DB, offset, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}
func (a *App) GetPost(w http.ResponseWriter, r *http.Request) {

}

func (a *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p models.Post

	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.Timestamp = time.Now().Format("2006-01-02T15:04:05Z")
	p.Post_ID = xid.New().String()
	if err := p.CreatePost(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, p)
}
func (a *App) DeletePost(w http.ResponseWriter, r *http.Request) {

}
