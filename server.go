package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "modernc.org/sqlite"
)

type Goods struct {
	Name        string
	Description string
	Price       int
	Weight      string
	ImagePath   string
}
type Examples struct {
	Name      string
	Price     int
	ImagePath string
}

func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite", "marmelad.db")
	if err != nil {
		fmt.Println("Ошибка подключения к БД:", err)
		return nil
	}
	return db
}

func selectAllGoods(db *sql.DB) []Goods {
	rows, err := db.Query("SELECT *FROM marmelad")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	sweets := []Goods{}
	for rows.Next() {
		a := Goods{}
		err := rows.Scan(&a.Name, &a.Description, &a.Price, &a.Weight, &a.ImagePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sweets = append(sweets, a)
	}
	return sweets
}

func selectAllExamples(db *sql.DB) []Examples {
	rows, err := db.Query("SELECT name, price, imagepath FROM marmelad LIMIT 4")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	fourSweets := []Examples{}
	for rows.Next() {
		a := Examples{}
		err := rows.Scan(&a.Name, &a.Price, &a.ImagePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fourSweets = append(fourSweets, a)
	}
	return fourSweets
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	db := dbConnect()
	defer db.Close()

	if r.Method == "GET" {
		fourSweets := selectAllExamples(db)
		tmpl, err := template.ParseFiles("templates/1page.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, fourSweets)
		if err != nil {
			fmt.Println("Ошибка выполнения шаблона:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	db := dbConnect()
	defer db.Close()
	if r.Method == "GET" {
		sweets := selectAllGoods(db)
		tmpl, err := template.ParseFiles("templates/2page.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, sweets)
		if err != nil {
			fmt.Println("Ошибка выполнения шаблона:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/3page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/4page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func QuestionsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/5page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/catalog", CatalogHandler)
	http.HandleFunc("/question_answer", QuestionsHandler)
	http.HandleFunc("/about-us", AboutHandler)
	http.HandleFunc("/contacts", ContactsHandler)
	fmt.Println("Server is running...")
	http.ListenAndServe(":8080", nil)

}
