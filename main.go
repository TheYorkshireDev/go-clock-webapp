package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  string
}

// Load all templates all on startup
var templates = template.Must(template.ParseGlob("templates/*"))

func newRouter() *mux.Router {
	dir := "assets"

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/time", timeHandler).Methods("GET")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir)))).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(pageNotFoundHandler)

	return router
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", &Page{Title: "Index", Body: "This is the HomePage"})
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "time", &Page{Title: "Time", Body: getTimeNow()})
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "404", &Page{Title: "Page Not Found", Body: "This is not the page your looking for"})
}

func getTimeNow() string {
	timeNow := time.Now()
	formattedTime := timeNow.Format("Mon Jan 02 15:04:05 MST 2006")

	return formattedTime
}

func main() {
	router := newRouter()

	http.ListenAndServe(":8080", router)
}
