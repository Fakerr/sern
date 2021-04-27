package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fakerr/sern/http/session"
	"github.com/Fakerr/sern/persist"
)

// return the user's repositories list
func GetRepositoriesList(w http.ResponseWriter, r *http.Request) {

	sess := session.Instance(r)
	user := sess.Values["login"].(string)

	repos, err := persist.GetRepositoriesByOwner(user)
	if err != nil {
		log.Printf("ERRO: persist.GetRepositoriesByOwner() failed with '%s'\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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
