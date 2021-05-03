package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Fakerr/sern/http/handlers"
	"github.com/Fakerr/sern/http/handlers/api"
	"github.com/Fakerr/sern/http/handlers/authentication"
	"github.com/Fakerr/sern/http/handlers/hooks"
	"github.com/Fakerr/sern/http/middleware"
	"github.com/Fakerr/sern/http/session"
	"github.com/Fakerr/sern/persist"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const (
	PUBLIC_DIR = "/public/"
)

// Force SSL redirection
func ForceSsl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "production" {
			if r.Header.Get("x-forwarded-proto") != "https" {
				sslUrl := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslUrl, http.StatusTemporaryRedirect)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	// Init postgres connection
	persist.InitConnection()

	// Init database tables
	persist.InitTables()

	// Init Redis connection pool
	persist.InitRedis()

	// Configure the session cookie store
	session.Configure()

	r := mux.NewRouter()

	// Middlewares
	r.Use(middleware.LogRequest)

	r.PathPrefix(PUBLIC_DIR).Handler(http.StripPrefix(PUBLIC_DIR, http.FileServer(http.Dir("."+PUBLIC_DIR))))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/build/static"))))

	r.HandleFunc("/", handlers.MainHandler)
	r.HandleFunc("/login", authentication.LoginHandler)
	r.HandleFunc("/logout", authentication.LogoutHandler)
	r.HandleFunc("/github_oauth_cb", authentication.GithubCallbackHandler)
	r.HandleFunc("/github_webhook_cb", hooks.WebhookCallbackHandler)

	// API requiring authenitcation.
	r.Handle("/api/repos", alice.New(middleware.WithAuth).ThenFunc(api.GetRepositoriesList)).Methods("GET")

	// API not requiring authenitcation
	r.Handle("/api/{owner}/{repo}/queue", alice.New().ThenFunc(api.GetQueue)).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		// default port
		port = "8089"
	}

	log.Print("INFO: Started running on http://127.0.0.1:" + port + "\n")
	log.Fatal(http.ListenAndServe(":"+port, ForceSsl(r)))
}
