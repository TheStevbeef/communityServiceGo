package app

import (
	"database/sql"
	"encoding/json"
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

	// statement, err = a.DB.Prepare("CREATE TABLE IF NOT EXISTS user (user_id TEXT PRIMARY KEY, name TEXT, image_url TEXT, post_id TEXT, FOREIGN KEY(post_id) REFERENCES post(post_id))")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// statement.Exec()
	statement, err := a.DB.Prepare("CREATE TABLE IF NOT EXISTS post (post_id TEXT PRIMARY KEY, timestamp TEXT, message TEXT, user_id)") //There should be this relationship: FOREIGN KEY(user_id) REFERENCES user(user_id)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatal(err)
	}
	statement, err = a.DB.Prepare("CREATE TABLE IF NOT EXISTS media (content_type TEXT, url TEXT, post_id TEXT, FOREIGN KEY(post_id) REFERENCES post(post_id))")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := statement.Exec(); err != nil {
		log.Fatal(err)
	}
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
	//a.Router.Use(authorization.RequireTokenAuthorization)
}

func (a *App) GetPosts(w http.ResponseWriter, r *http.Request) {
	var limit int
	var offset int
	limit = 10
	offset = 10
	if r.URL.Query().Get("offset") != "" &&
		r.URL.Query().Get("offset") != "" {
		var err error
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if limit > 10 || limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	posts, err := models.GetPosts(a.DB, offset, limit)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//Dummy functions to get the Users
	posts, err = models.GetUsers(posts)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, posts)
}

func (a *App) GetPost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	//Get Post
	p := models.Post{Post_ID: id}
	if err := p.GetPost(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// Add User to Post
	if err := p.GetUser(); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, p)
}

func (a *App) CreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p models.Post
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if !IsPayloadValid(p) {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	p.Timestamp = time.Now().Format("2006-01-02T15:04:05Z")
	p.Post_ID = xid.New().String()
	if err := p.CreatePost(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"result": "Successfully created"})
}

func IsPayloadValid(post models.Post) bool {
	if post.User.User_ID != "" && post.User.Name != "" &&
		(post.Message != "" || (post.Media.Content_type != "" && post.Media.Url != "")) &&
		((post.Media.Content_type != "" && post.Media.Url != "") || (post.Media.Content_type == "" && post.Media.Url == "")) &&
		(len(post.Message) <= 140) {
		return true
	} else {
		return false
	}
}

func (a *App) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	p := models.Post{Post_ID: id}
	if err := p.DeletePost(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "successfully deleted"})

}
