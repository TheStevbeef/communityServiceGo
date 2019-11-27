package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Post_ID   string `json:"id"`
	Timestamp string `json:"timestamp"`
	User      User   `json:"user"`
	Message   string `json:"message"`
	Media     Media  `json:"media"`
}

type User struct {
	User_ID   string `json:"id"`
	Name      string `json:"name"`
	Image_url string `json:"image_url"`
}

type Media struct {
	Content_type string `json:"content_type"`
	Url          string `json:"url"`
}

func (p *Post) GetPost(db *sql.DB) error {
	post_statement := fmt.Sprintf("SELECT post_id,timestamp,message FROM post WHERE post_id='%s'", p.Post_ID)
	user_statement := fmt.Sprintf("SELECT user_id,name,image_url FROM user WHERE post_id='%s'", p.Post_ID)
	media_statement := fmt.Sprintf("SELECT content_type,url FROM media WHERE post_id='%s'", p.Post_ID)
	row := db.QueryRow(post_statement)
	if err := row.Scan(&p.Post_ID, &p.Timestamp, &p.Message); err != nil {
		return err
	}
	row = db.QueryRow(user_statement)
	if err := row.Scan(&p.User.User_ID, &p.User.Name, &p.User.Image_url); err != nil {
		return err
	}
	row = db.QueryRow(media_statement)
	if err := row.Scan(&p.Media.Content_type, &p.Media.Url); err != nil {
		return err
	}
	return nil
}

func (p *Post) DeletePost(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM post WHERE post_id='%s'", p.Post_ID)
	if _, err := db.Exec(statement); err != nil {
		return err
	}
	statement = fmt.Sprintf("DELETE FROM media WHERE post_id='%s'", p.Post_ID)
	if _, err := db.Exec(statement); err != nil {
		return err
	}
	return nil

}

func (p *Post) CreatePost(db *sql.DB) error {
	//TODO: User ID if it is twice then it should also work properly

	// Insert into post
	statement := fmt.Sprintf("INSERT INTO post(post_id, timestamp, message) VALUES('%s','%s','%s')", p.Post_ID, p.Timestamp, p.Message)
	if _, err := db.Exec(statement); err != nil {
		return err
	}
	//Insert into user
	statement = fmt.Sprintf("INSERT INTO user(user_id, name, image_url, post_id) VALUES('%s','%s','%s','%s')", p.User.User_ID, p.User.Name, p.User.Image_url, p.Post_ID)
	if _, err := db.Exec(statement); err != nil {
		return err
	}
	//Insert into media
	statement = fmt.Sprintf("INSERT INTO media(content_type, url, post_id) VALUES('%s','%s', '%s')", p.Media.Content_type, p.Media.Url, p.Post_ID)
	if _, err := db.Exec(statement); err != nil {
		return err
	}
	return nil
}

func GetPosts(db *sql.DB, offset, limit int) ([]Post, error) {
	statement := fmt.Sprintf("SELECT post_id FROM post LIMIT %d OFFSET %d", limit, offset)
	println("asdf")
	post_ids, err := db.Query(statement)
	if err != nil {
		return nil, err
	}

	defer post_ids.Close()
	posts := []Post{}
	for post_ids.Next() {
		var post_id string
		if err := post_ids.Scan(&post_id); err != nil {
			log.Fatal(err)
		}
		post_statement := fmt.Sprintf("SELECT post_id,timestamp,message FROM post WHERE post_id='%s'", post_id)
		user_statement := fmt.Sprintf("SELECT user_id,name,image_url FROM user WHERE post_id='%s'", post_id)
		media_statement := fmt.Sprintf("SELECT content_type,url FROM media WHERE post_id='%s'", post_id)
		var p Post
		row := db.QueryRow(post_statement)
		if err := row.Scan(&p.Post_ID, &p.Timestamp, &p.Message); err != nil {
			return nil, err
		}
		row = db.QueryRow(user_statement)
		if err := row.Scan(&p.User.User_ID, &p.User.Name, &p.User.Image_url); err != nil {
			return nil, err
		}
		row = db.QueryRow(media_statement)
		if err := row.Scan(&p.Media.Content_type, &p.Media.Url); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
