package handlers

import "net/http"

// Define All your handlers here
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home from handler!"))
}