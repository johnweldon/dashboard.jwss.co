package main

import (
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/unrolled/secure"

	"github.com/johnweldon/dashboard.jwss.co/area/test"
)

func main() {
	router := mux.NewRouter()
	test.RegisterArea(router.PathPrefix("/test").Subrouter())

	secureMiddleware := secure.New(secure.Options{
		SSLRedirect:        true,
		SSLHost:            "secure.jwss.co/dashboard",
		SSLProxyHeaders:    map[string]string{"X-Forwarded-Proto": "https"},
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IsDevelopment:      false,
	})

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)
	n.Run(listenEndpoint())
}

func listenEndpoint() string {
	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	return "0.0.0.0:" + port
}
