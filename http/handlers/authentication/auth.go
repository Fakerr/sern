package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/http/session"
)

// login handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	sess := session.Instance(r)
	sess.Values["state"] = state
	sess.Save(r, w)

	url := config.OauthConf.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	session.Empty(sess)

	sess.Options.MaxAge = -1
	sess.Save(r, w)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
