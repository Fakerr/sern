package api

import (
	"encoding/json"
	"net/http"

	"github.com/Fakerr/sern/cors/runner"

	"github.com/gorilla/mux"
)

// Return the merge queue for a specific repository
func GetQueue(w http.ResponseWriter, r *http.Request) {

	owner := mux.Vars(r)["owner"]
	repo := mux.Vars(r)["repo"]

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

	// return only the PR's number
	var parsedPRs []int
	for _, item := range queue {
		parsedPRs = append(parsedPRs, item.Number)
	}

	js, _ := json.Marshal(parsedPRs)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
