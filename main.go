package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
	temp, _ := template.ParseFiles("templates/register.html")
	name := r.FormValue("name")
	pass := r.FormValue("password")
	email := r.FormValue("email")
	if name == "" || pass == "" || email == "" {
		w.WriteHeader(http.StatusBadRequest)
		temp.Execute(w, "please fill the form!")
		return
	}
	post := `INSERT INTO User (email,username,password_hash)
	VALUES (?, ?, ?)`
	_, err := db.Exec(post, email, name, pass)
	if err != nil {
		temp.Execute(w, "email or name already exists")
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	temp, _ := template.ParseFiles("templates/login.html")
	name := r.FormValue("name")
	pass := r.FormValue("password")

	if name == "" || pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		temp.Execute(w, "please fill the form!")
		return
	}
	post := `select password_hash,deja from user where username = ?`
	f := db.QueryRow(post, name)
	if f.Err() != nil {
		temp.Execute(w, "username not found !")
		return
	}
	var passw string
	var deja int
	f.Scan(&passw, &deja)
	if pass != passw {
		temp.Execute(w, "incorrect password !")
		return
	}
	c := http.Cookie{
		Name:     "username",
		Value:    name,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}
	if deja == 0 {
		db.Exec(`UPDATE User SET deja = ? WHERE username = ?`, 1, name)
	} else {
		temp.Execute(w, "chi 9lwa")
		return
	}
	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func register(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/register.html")
	temp.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/login.html")
	temp.Execute(w, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("templates/home.html")
	temp.Execute(w, "name.Value")
}

func logout(w http.ResponseWriter, r *http.Request) {
	name, _ := r.Cookie("username")
	c := http.Cookie{
		Name:   "username",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
	post := `select deja from user where username = ?`
	f := db.QueryRow(post, name.Value)
	var deja int
	f.Scan(&deja)
	fmt.Println(deja)
	if deja == 1 {
		db.Exec(`UPDATE User SET deja = ? WHERE username = ?`, 0, name.Value)
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS User (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		deja INT DEFAULT 0
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/handleRegister", handleRegister)
	http.HandleFunc("/handleLogin", handleLogin)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/", home)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}
