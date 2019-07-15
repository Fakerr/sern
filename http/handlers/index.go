package handlers

import (
	"html/template"
	"net/http"

	"github.com/Fakerr/sern/http/session"
)

var templates = template.Must(template.ParseFiles("public/home.html", "ui/build/index.html"))

// Hnadler for the main route
func MainHandler(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	// Check whether this is a new session or the user is not authenticated.
	if sess.IsNew == true || sess.Values["id"] == nil {
		templates.ExecuteTemplate(w, "home.html", nil)
		return
	}
	templates.ExecuteTemplate(w, "index.html", nil)
}
