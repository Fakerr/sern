package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fakerr/sern/cors/client"
	"github.com/google/go-github/github"

	"github.com/gorilla/mux"
)

type repository struct {
	id    string
	owner string
	// The first time the user enable a repo, it should be presisted with its config.
	// Later the user can disable the repo but its config is still saved.
	enabled boolean
}

var enabledRepositories []repository

// Fetch and return user's repositories
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

// PRs for enabled repos will be queued in the merge queue
func EnableRepository(w http.ResponseWriter, r *http.Request) {
	sess := session.Instance(r)

	//client := client.FromToken(r.Context(), sess.Values["accessToken"])

	repo := repository{id: r.FormValue["repoID"], owner: sess.Values["id"], enabled: true}

	enabledRepositories := append(enabledRepositories, repo)
}
