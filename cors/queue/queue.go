package queue

import (
	"log"
)

// A map with a all enabled repositories and their current queue. Should be presisted (redis).
var ReposQueue map[string][]*PullRequest

// Initilaize database queue
func SetQueue() {
	ReposQueue = make(map[string][]*PullRequest)
}

type PullRequest struct {
	Id     int
	Status string
}

// If the merge queue doesn't exist, a new one is created for the concerned repository.
func AddToQueue(repo string, pr *PullRequest) error {

	log.Printf("ReposQueue %s\n", ReposQueue)

	ReposQueue[repo] = append(ReposQueue[repo], pr)

	log.Printf("ReposQueue %s\n", ReposQueue)

	return nil
}

// next will check whether or not a PR is already being processed and if not will start
// processing the next PR in the merge queue.
func next(repo string) {

}
