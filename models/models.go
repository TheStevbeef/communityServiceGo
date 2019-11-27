package models

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type PostWithoutIDAndTime struct {
	User    User   `json:"user"`
	Message string `json:"message"`
	Media   Media  `json:"media"`
}

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
	// statement := fmt.Sprintf("SELECT name, age FROM users WHERE id=%d", u.ID)
	// return db.QueryRow(statement).Scan(&u.Name, &u.Age)
	return nil
}

func (p *Post) DeletePost(db *sql.DB) error {
	// statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	// _, err := db.Exec(statement)
	// return err
	return nil
}

func (p *Post) CreatePost(db *sql.DB) error {
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
	// statement := fmt.Sprintf("SELECT id, name, age FROM users LIMIT %d OFFSET %d", count, start)
	// rows, err := db.Query(statement)

	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()

	// users := []user{}

	// for rows.Next() {
	// 	var u user
	// 	if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, u)
	// }

	// return users, nil
	return nil, nil
}
