package test

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterArea(router *mux.Router) {
	router.HandleFunc("/", bareTest)
	router.HandleFunc("/{name}", nameTest)
}

func bareTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Bare Necessities")
}

func nameTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	fmt.Fprintf(w, "Hello %s", name)
}
