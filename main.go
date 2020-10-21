package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Port is the port the webserver runs on.
const (
	Port = ":8080"
)

// Global variable to pass quote through to html template
var quote string

// Random number generator
func quotenum() int {
	rand.Seed(time.Now().UnixNano())
	var returnint int = rand.Intn(10-1) + 1
	return returnint
}

// Html template generation
func serveStatic(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("quote.html")
	if err != nil {
		fmt.Println(err)
	}

	items := struct {
		Quote string
	}{
		Quote: postgresQuery(),
	}
	t.Execute(w, items)
}

// Run a simple query against postgres to return one row
func postgresQuery() string {
	psqlconn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	err = db.QueryRow("SELECT quote FROM quotes where id=$1", quotenum()).Scan(&quote)
	CheckError(err)
	return quote

}

// CheckError is a basic error handler
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", serveStatic)
	http.ListenAndServe(Port, nil)
}
