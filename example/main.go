package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/taironas/auth"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/token", tokenHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./index.html")
	t.Execute(w, nil)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	IDToken := r.FormValue("id_token")

	if user, err := auth.Google(IDToken); err != nil {
		fmt.Fprintf(w, "error: %v", err)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}
