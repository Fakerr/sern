package middleware

import (
	"log"
	"net/http"

	"github.com/Fakerr/sern/http/session"
)

// Do not allow anonymous users to access the ressource
func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess := session.Instance(r)

		// If users are not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			//http.Redirect(w, r, "/", http.StatusUnauthorized)
			http.Error(w, "401 not authorized", http.StatusUnauthorized)
			return
		}

		log.Printf("INFO: Authenticated user %s\n", sess.Values["login"])

		next.ServeHTTP(w, r)
	})
}
