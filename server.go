package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "modernc.org/sqlite"
)

var httpRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total numbers of HTTP Requests",
	},
	[]string{"path"},
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
type Feedbacks struct {
	Name     string
	Email    string
	Comments string
}

func dbConnect1() *sql.DB {
	db, err := sql.Open("sqlite", "marmelad.db")
	if err != nil {
		fmt.Println("Ошибка подключения к БД1:", err)
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
	httpRequestTotal.WithLabelValues(r.URL.Path).Inc()

	db := dbConnect1()
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

	httpRequestTotal.WithLabelValues(r.URL.Path).Inc()

	db := dbConnect1()
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

	httpRequestTotal.WithLabelValues(r.URL.Path).Inc()

	tmpl, err := template.ParseFiles("templates/3page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func ContactsHandler(w http.ResponseWriter, r *http.Request) {

	httpRequestTotal.WithLabelValues(r.URL.Path).Inc()

	tmpl, err := template.ParseFiles("templates/4page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func QuestionsHandler(w http.ResponseWriter, r *http.Request) {

	httpRequestTotal.WithLabelValues(r.URL.Path).Inc()

	tmpl, err := template.ParseFiles("templates/5page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func FeedbackHandler(w http.ResponseWriter, r *http.Request) {
	httpRequestTotal.WithLabelValues(r.URL.Path).Inc()

	if r.Method == "POST" {
		db := dbConnect1()
		defer db.Close()
		var person Feedbacks
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&person)
		if err != nil {
			fmt.Println(err)
			return
		}

		if person.Name == "" || person.Email == "" || person.Comments == "" {
			fmt.Println(err)
			return
		}

		_, err = db.Exec("INSERT INTO feedbacks (name, email, comments) VALUES (?, ?, ?)", person.Name, person.Email, person.Comments)
		if err != nil {
			fmt.Println("Ошибка вставки в БД:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Feedback saved successfully",
		})

	} else {
		tmpl, err := template.ParseFiles("templates/6page.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func main() {

	prometheus.MustRegister(httpRequestTotal)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/catalog", CatalogHandler)
	http.HandleFunc("/question_answer", QuestionsHandler)
	http.HandleFunc("/about-us", AboutHandler)
	http.HandleFunc("/contacts", ContactsHandler)
	http.HandleFunc("/feedback", FeedbackHandler)
	fmt.Println("Server is running...")
	http.ListenAndServe(":8080", nil)

}
