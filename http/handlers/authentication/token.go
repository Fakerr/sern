package authentication

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/cors/client"
	"github.com/Fakerr/sern/server/session"
)

// /github_oauth_cb. Called by github after authorization is granted
// If the user is successfully authenticated, create and save the user's session information
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	sess := session.Instance(r)

	state := r.FormValue("state")
	if state != sess.Values["state"] {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", sess.Values["state"], state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := config.OauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//client := github.NewClient(oauthConf.Client(oauth2.NoContext, token))
	client := client.FromToken(oauth2.NoContext, token.AccessToken)

	user, _, err := client.Users.Get(context.Background(), "")

	if err != nil {
		log.Printf("client.Users.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sess.Values["id"] = user.ID
	sess.Values["login"] = user.Login
	sess.Values["userName"] = user.Name
	sess.Values["accessToken"] = token.AccessToken
	sess.Save(r, w)

	log.Printf("Logged in as GitHub user: %s\n", *user.Login)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
