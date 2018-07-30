package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/Fakerr/sern/server/session"
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

		// Extend context with the user's token
		if token, ok := sess.Values["accessToken"]; ok {
			log.Printf("Authenticated user %s\n", sess.Values["login"])
			ctx := context.WithValue(r.Context(), "token", token)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
