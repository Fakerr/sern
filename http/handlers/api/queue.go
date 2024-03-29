package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fakerr/sern/core/runner"
	"github.com/Fakerr/sern/http/session"
	"github.com/Fakerr/sern/persist"

	"github.com/gorilla/mux"
)

// Return the merge queue for a specific repository
func GetQueue(w http.ResponseWriter, r *http.Request) {

	owner := mux.Vars(r)["owner"]
	repo := mux.Vars(r)["repo"]

	// Check whether or not the repository exist
	exist, err := persist.RepositoryExists(owner + "/" + repo)
	if err != nil {
		log.Printf("ERRO: persist.RepositoryExists() failed with '%s'\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !exist {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	isPrivate, err := persist.IsPrivate(owner + "/" + repo)
	if err != nil {
		log.Printf("ERRO: persist.IsPrivate() failed with '%s'\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isPrivate == true {
		if auth := session.IsAuthenticated(r); auth == false {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		sess := session.Instance(r)
		login := sess.Values["login"].(string)
		if login != owner {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}
	}

	// Get the Repo's runner
	runner := runner.GetSoftRunner(owner, repo)

	// If there is no runner for the repo (no items in the queue), return null as a response
	// Not sure if this is the best way to do it
	if runner == nil {
		js, _ := json.Marshal(nil)

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	queue := runner.Queue.QueueItems
	active := runner.Active.Number
	status := runner.Status

	// return only the PR's number
	var parsedPRs []int
	for _, item := range queue {
		parsedPRs = append(parsedPRs, item.Number)
	}

	res := &QueueResponse{
		Queue:  parsedPRs,
		Active: active,
		Status: status,
	}

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
