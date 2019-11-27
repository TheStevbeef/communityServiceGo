package main

import (
	"os"

	_ "github.com/TheStevbeef/communityServiceGo/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// os.Setenv("DB_PATH", ".\\Community.db")
	// fmt.Println(os.Environ())
	a := App{}
	a.Initialize(os.Getenv("DB_PATH"))

	a.Run(":8080")
}
