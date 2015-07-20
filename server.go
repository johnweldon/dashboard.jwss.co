package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/unrolled/secure"

	"github.com/johnweldon/dashboard.jwss.co/area/test"
)

func init() {
	initializeAreas()
	initializeCompression()
	initializeSecureMiddleware()
	initializeSessions()
	initializeStore()
}

var (
	store            *sessions.CookieStore
	secureMiddleware *secure.Secure
	gorillaRouter    *mux.Router

	secureHandler      negroni.Handler
	sessionHandler     negroni.Handler
	compressionHandler negroni.Handler
)

func main() {
	n := negroni.Classic()

	// compression first
	n.Use(compressionHandler)
	n.Use(secureHandler)
	n.Use(sessionHandler)

	n.UseHandler(gorillaRouter)

	n.Run(listenEndpoint())
}

func initializeSessions() {
	sessionHandler = negroni.HandlerFunc(sessionFunc)
}

func initializeStore() {
	store = sessions.NewCookieStore(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   120,
		HttpOnly: true,
		Secure:   true,
	}
}

func initializeSecureMiddleware() {
	secureMiddleware = secure.New(secure.Options{
		SSLRedirect:        true,
		SSLHost:            "secure.jwss.co/dashboard",
		SSLProxyHeaders:    map[string]string{"X-Forwarded-Proto": "https"},
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IsDevelopment:      false,
	})
	secureHandler = negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext)
}

func initializeCompression() {
	compressionHandler = gzip.Gzip(gzip.DefaultCompression)
}

func initializeAreas() {
	gorillaRouter = mux.NewRouter().StrictSlash(true)
	gorillaRouter.HandleFunc("/dashboard/login", loginFunc).Methods("POST")
	gorillaRouter.HandleFunc("/dashboard/logout", logoutFunc).Methods("POST")
	test.RegisterArea(gorillaRouter.PathPrefix("/test/").Subrouter())
}

func sessionFunc(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, _ := store.Get(r, "_jwss.dash")
	session.Save(r, w)

	if r.URL.Path == "/dashboard/" || session.Values["authenticated"] == "true" {
		if next != nil {
			next(w, r)
		}
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}

}

func loginFunc(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "_jwss.dash")
	user := r.FormValue("name")
	pass := r.FormValue("password")
	redirect := "/dashboard"
	session.Values["authenticated"] = "false"
	if user != "" && pass != "" {
		if user == "john" && pass == "john" {
			session.Values["authenticated"] = "true"
			redirect = "/dashboard/test"
		}
	}
	session.Save(r, w)
	http.Redirect(w, r, redirect, http.StatusFound)
}
func logoutFunc(w http.ResponseWriter, r *http.Request) {}

func listenEndpoint() string {
	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	return "0.0.0.0:" + port
}
