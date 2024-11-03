package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func handleRegister(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/index.html")
	name := r.FormValue("name")
	pass := r.FormValue("password")
	email := r.FormValue("email")

	if name == "" || pass == "" {
		temp.Execute(w, "error")
	}
	post := `INSERT INTO User (email,username,password_hash)
	VALUES (?, ?, ?)`
	_, err := db.Exec(post, email, name, name)
	if err != nil {
		temp.Execute(w, "error")
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/index.html")
	name := r.FormValue("name")
	pass := r.FormValue("password")

	if name == "" || pass == "" {
		temp.Execute(w, "error")
	}
	post := `select password_hash from user where username = ?`
	f := db.QueryRow(post, name)
	passw := ""
	f.Scan(&passw)
	fmt.Println(passw)
}

func register(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/index.html")
	temp.Execute(w, nil)
}
func login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/login.html")
	temp.Execute(w, nil)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err) // Log and exit if the connection fails
	}
	defer db.Close() // Close the database connection when main exits

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS User (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		log.Fatal(err) // Log and exit if the table creation fails
	}
	http.HandleFunc("/handleRegister", handleRegister)
	http.HandleFunc("/handleLogin", handleLogin)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.ListenAndServe(":8080", nil)
}
