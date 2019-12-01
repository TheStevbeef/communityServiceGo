The Community Service is a Restful API where you can read create and delete posts.

If you want to run this Service you need
- SQLite3 on your machine
- go get "github.com/mattn/go-sqlite3"
- go get "github.com/gorilla/mux"
- go get "github.com/rs/xid"
- go get "github.com/joho/godotenv"

The rough architecture is based on this github repository
https://github.com/kelvins/GoApiTutorial

The main file is for the initialization of the service.

Next there is the app.go file. This file handles the REST-Requests.
For the Routing I used the gorilla mux library, because it is very lightweight.
After a successfull request the request will be send to the model.go file.

There is happening all the Database commands in sqlite3.

In The Service is one environment variable. This is the Path, where you want to safe your database.
For example if you want to safe your database file in the root folder with the name "community.db", then you have to set an environment variable "DB_PATH" with the value "./community.db"

To test at the one hand the REST-API and on the other hand all the database-commands. There is a main_test.go file.