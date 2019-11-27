package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type post struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	User      user   `json:"user"`
	Message   string `json:"message"`
	Media     media  `json:"media"`
}

type user struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Image_url string `json:"image_url"`
}

type media struct {
	Content_type string `json:"content_type"`
	Url          string `json:"url"`
}
type posts struct {
	Posts []post `json:"post"`
}

func (p *post) getPost(db *sql.DB) error {
	// statement := fmt.Sprintf("SELECT name, age FROM users WHERE id=%d", u.ID)
	// return db.QueryRow(statement).Scan(&u.Name, &u.Age)
	return nil
}

func (p *post) deletePost(db *sql.DB) error {
	// statement := fmt.Sprintf("DELETE FROM users WHERE id=%d", u.ID)
	// _, err := db.Exec(statement)
	// return err
	return nil
}

func (p *post) createPost(db *sql.DB) error {
	// statement := fmt.Sprintf("INSERT INTO users(name, age) VALUES('%s', %d)", u.Name, u.Age)
	// _, err := db.Exec(statement)

	// if err != nil {
	// 	return err
	// }

	// err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	// if err != nil {
	// 	return err
	// }

	// return nil
	return nil
}

func getPosts(db *sql.DB, start, count int) ([]post, error) {
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

// database, _ := sql.Open("sqlite3", "./nraboy.db")
// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
// statement.Exec()
// statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
// statement.Exec("Nic", "Raboy")
// rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
// var id int
// var firstname string
// var lastname string
// for rows.Next() {
// 	rows.Scan(&id, &firstname, &lastname)
// 	fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
// }
