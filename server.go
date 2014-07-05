package main

import (
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"

	"github.com/johnweldon/dashboard.jwss.co/area/test"
)

func main() {
	router := mux.NewRouter()
	test.RegisterArea(router.PathPrefix("/test").Subrouter())

	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)
	n.Run(listenEndpoint())
}

func listenEndpoint() string {
	port := "9191"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	return "0.0.0.0:" + port
}
