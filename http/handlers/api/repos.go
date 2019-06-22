package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/core/client"
	"github.com/Fakerr/sern/http/session"
	"github.com/Fakerr/sern/persist"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

// return the user's repositories list
func GetRepositoriesList(w http.ResponseWriter, r *http.Request) {

	sess := session.Instance(r)
	user := sess.Values["login"].(string)

	repos := persist.GetRepositoriesByOwner(user)

	// Extract the repos' name
	// using make([]string, 0) instead of var parsedRepos []string to handle the case of empty reposonse
	// explanation: https://www.danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	parsedRepos := make([]string, 0)
	for _, repo := range repos {
		parsedRepos = append(parsedRepos, repo.FullName)
	}

	js, _ := json.Marshal(parsedRepos)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// TODO remove
func RepositoriesList(w http.ResponseWriter, r *http.Request) {

	client := github.NewClient(nil)

	user := mux.Vars(r)["user"]

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	// fetch all user's repositories
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(r.Context(), user, opt)
		if err != nil {
			log.Printf("client.Repositories.List() failed with '%s'\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// return only the repo's name
	var parsedRepos []string
	for _, repo := range allRepos {
		parsedRepos = append(parsedRepos, *repo.Name)
	}

	js, _ := json.Marshal(parsedRepos)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Pull requests for enabled repos will be queued before they get merged.
// When the user enable a repository, a webhook will be created and will listen
// to any Pull request event.
// The repositry config will be persisted in the database.
func EnableRepository(w http.ResponseWriter, r *http.Request) {

	log.Print("Processing EnableRepository ...\n")

	sess := session.Instance(r)

	token := sess.Values["accessToken"].(string)
	userLogin := sess.Values["login"].(string) // Maybe this one should be reconsidered...
	repo := r.FormValue("repoID")
	fullRepo := userLogin + "/" + repo // Should be in this form: owner/repo

	// Create the repository and presist it if it doesn't exist in the db.
	repository := &persist.Repository{
		FullName: fullRepo,
		Owner:    userLogin,
		Private:  false,
	}

	// Enable repository (Persist in the db).
	err := persist.AddRepository(repository)
	if err != nil {
		log.Printf("persist.addRepository() failed with '%s'\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create a webhook for this repository
	// Get a github client using the user's token
	client := client.FromToken(r.Context(), token)

	// Get a default hook config
	hook := config.GetHookConfig(userLogin)

	// Normally, this should not return a 'webhook already exists' error since each webhook
	// will be deleted once the user disable a repository.
	_, _, err = client.Repositories.CreateHook(r.Context(), userLogin, repo, hook)
	if err != nil {
		log.Printf("client.Repositories.CreateHook() failed with '%s'\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
