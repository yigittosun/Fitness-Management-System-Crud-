package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Member struct {
	Id   int
	Fullname string
	Height int
	Weight int
	Age int
	Program int
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "deneme"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var selDB, err = db.Query("SELECT * FROM members ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	obje := Member{}
	liste := []Member{}
	for selDB.Next() {
		var id,height,weight,age,program int
		var name string

		err = selDB.Scan(&id, &name, &height, &weight, &age, &program)
		if err != nil {
			panic(err.Error())
		}
		obje.Id = id
		obje.Fullname = name
		obje.Height = height
		obje.Weight = weight
		obje.Age=age
		obje.Program = program
		liste = append(liste, obje)
	}
	tmpl.ExecuteTemplate(w, "2.html", liste)
	defer db.Close()
}


func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		age := r.FormValue("age")
		height := r.FormValue("height")
		weight := r.FormValue("weight")
		program := r.FormValue("program")
		insForm, err := db.Prepare("INSERT INTO members(Fullname, Age, Height, Weight, Program) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, age, height, weight, program)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}


func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	obje := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM members WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(obje)
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
func About(w http.ResponseWriter, r *http.Request){
	liste := []Member{}
	tmpl.ExecuteTemplate(w, "About.html", liste)
}
func Blog(w http.ResponseWriter, r *http.Request){
	liste := []Member{}
	tmpl.ExecuteTemplate(w, "Blog.html", liste)
}

func Equipment(w http.ResponseWriter, r *http.Request){
	liste := []Member{}
	tmpl.ExecuteTemplate(w, "Equipment.html", liste)
}
func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/about", About)
	http.HandleFunc("/Blog", Blog)
	http.HandleFunc("/Equipment", Equipment)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
