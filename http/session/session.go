package session

import (
	"net/http"

	"github.com/Fakerr/sern/config"

	"github.com/gorilla/sessions"
)

var (
	// Store is the cookie store
	Store *sessions.CookieStore
	// Name is the session name
	Name string
)

// Configure the session cookie store
func Configure() {
	secretKey := config.SessionSecretKey
	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   60 * 60 * 24,
	}
	Name = "_uuid"
}

// Instance returns a new session, never returns an error
func Instance(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, Name)
	return session
}

// Empty deletes all the current session values
func Empty(sess *sessions.Session) {
	// Clear out all stored values in the cookie
	for k := range sess.Values {
		delete(sess.Values, k)
	}
}

// Return true if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	sess := Instance(r)
	if sess.Values["id"] == nil {
		return false
	}
	return true
}
