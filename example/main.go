package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taironas/auth"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/api/token", tokenHandler)
	http.ListenAndServe(":8080", nil)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	IDToken := r.FormValue("id_token")

	if user, err := auth.Google(IDToken); err != nil {
		fmt.Fprintf(w, "error: %v", err)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}
