package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/TheStevbeef/communityServiceGo/app"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	a := app.App{}
	a.Initialize(os.Getenv("DB_PATH"))
	a.Run(":8080")
}
