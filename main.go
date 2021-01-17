package main

import (
  "fmt"
  "net/http"
  //"text/template"
  "html/template"
)

// process form inputs with go templates
/*
https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/04.1.html
*/

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello All</h1>")
}

func search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: " + r.Method)

	if r.Method == "GET" {
		t,_ := template.ParseFiles("search.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		// process data
		fmt.Println("Search Query: ", r.Form["searchQuery"])
	}
}

func main() {
	fmt.Println("hello all")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/search", search)
	http.ListenAndServe(":8088", nil)
}