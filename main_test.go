// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/TheStevbeef/communityServiceGo/models"

	"github.com/TheStevbeef/communityServiceGo/app"
)

var a app.App

func TestMain(m *testing.M) {
	a = app.App{}
	a.Initialize(os.Getenv("DB_PATH_TEST"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(postTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(mediaTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE * FROM post")
	a.DB.Exec("DELETE * FROM media")
}

const postTableCreationQuery = `
CREATE TABLE IF NOT EXISTS post 
(post_id TEXT PRIMARY KEY, 
	timestamp TEXT, 
	message TEXT, 
	user_id)`
const mediaTableCreationQuery = `
CREATE TABLE IF NOT EXISTS media 
(content_type TEXT, 
	url TEXT, 
	post_id TEXT, FOREIGN KEY(post_id) REFERENCES post(post_id))`

func TestEmptyTable(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/posts", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentPost(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/posts/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

func TestCreatePost(t *testing.T) {
	clearTable()

	payload := []byte(`{
		"user": {
		  "id": "test_user",
		  "name": "Kevin",
		  "image_url": "https://i.pravatar.cc/150?u=abc12345"
		},
		"message": "123123",
		"media":{
			"content_type": "image/jpeg",
			"url": "https://picsum.photos/id/6/300"
		}
	}`)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
	var p models.Post
	row := a.DB.QueryRow("Select * FROM media")
	if err := row.Scan(&p.Media.Content_type, &p.Media.Url, &p.Post_ID); err != nil {
		t.Errorf("Something with the database went wrong")
	}
	row = a.DB.QueryRow("Select * FROM post")
	if err := row.Scan(&p.Post_ID, &p.Timestamp, &p.Message, &p.User.User_ID); err != nil {
		t.Errorf("Something with the database went wrong")
	}

	if p.User.User_ID != "test_user" {
		t.Errorf("Expected user name to be 'test_user'. Got '%v'", p.User.User_ID)
	}
	if p.Message != "123123" {
		t.Errorf("Expected Message to be '123123'. Got '%v'", p.Message)
	}
	if p.Media.Content_type != "image/jpeg" {
		t.Errorf("Expected content type to be 'image/jpeg'. Got '%v'", p.Media.Content_type)
	}
	if p.Media.Url != "https://picsum.photos/id/6/300" {
		t.Errorf("Expected media url to be 'https://picsum.photos/id/6/300'. Got '%v'", p.Media.Url)
	}

}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		postStatement := fmt.Sprintf("INSERT INTO post(post_id, timestamp,message, user_id) VALUES('%s', '%s', '%s','%s')", ("post_id_" + strconv.Itoa(i+1)), "2006-01-02T15:04:05Z", ("message_" + strconv.Itoa(i+1)), ("user_id_" + strconv.Itoa(i+1)))
		mediaStatement := fmt.Sprintf("INSERT INTO media(content_type, url, post_id) VALUES('%s', '%s', '%s')", ("content_type_" + strconv.Itoa(i+1)), ("url_" + strconv.Itoa(i+1)), ("post_id_" + strconv.Itoa(i+1)))
		a.DB.Exec(postStatement)
		a.DB.Exec(mediaStatement)
	}
}

func TestGetUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/posts/post_id_1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/posts/post_id_1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/posts/post_id_1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/posts/post_id_1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
