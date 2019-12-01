package models

import (
	_ "github.com/mattn/go-sqlite3"
)

// Dummy functions to get a completed post with the users.
func (p *Post) GetUser() error {
	p.User.Name = "Username" + p.User.User_ID
	p.User.Image_url = "URL" + p.User.User_ID
	return nil
}

func GetUsers(posts []Post) ([]Post, error) {

	for _, post := range posts {
		if err := post.GetUser(); err != nil {
			return nil, err
		}
	}
	return posts, nil
}
