package handlers

import (
	"html/template"
	"net/http"

	"github.com/Fakerr/sern/http/session"
)

var templates = template.Must(template.ParseFiles("static/logout.html", "static/index.html"))

// Hnadler for the main route
func MainHandler(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	if sess.IsNew == true {
		templates.ExecuteTemplate(w, "index.html", nil)
		return
	}
	templates.ExecuteTemplate(w, "logout.html", nil)
}