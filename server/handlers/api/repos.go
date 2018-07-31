package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
)

// fetch and return user's repositories
func GetUserRepos(w http.ResponseWriter, r *http.Request) {
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
