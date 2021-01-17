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
	//fmt.Fprintf(w, "<h1>Hello All</h1>")
	t,_ := template.ParseFiles("root.gtpl")
	t.Execute(w, nil)

}

func search(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: " + r.Method)
	t,_ := template.ParseFiles("search.gtpl")
	if r.Method == "GET" {
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		searchQuery := r.Form["searchQuery"]
		// process data
		if searchQuery != nil {
			fmt.Println("Search Query: ", searchQuery)
		} else {
			fmt.Println("Empty search query")
		}
		t.Execute(w, nil)

	}
}

func main() {
	fmt.Println("hello all")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/search", search)
	fmt.Println("listening on 0.0.0.0:8088")
	http.ListenAndServe(":8088", nil)
}