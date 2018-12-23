package queue

import (
	"context"
	"log"

	"github.com/Fakerr/sern/cors/actions"

	"github.com/google/go-github/github"
)

// A map with a all enabled repositories and their current queue. Should be presisted (redis).
var ReposQueue map[string][]*PullRequest

// Initilaize database queue
func SetQueue() {
	ReposQueue = make(map[string][]*PullRequest)
}

type PullRequest struct {
	Number         int
	Status         string
	HeadSHA        string
	MergeCommitSHA string
}

// If the merge queue doesn't exist, a new one is created for the concerned repository.
func AddToQueue(owner, repo string, pr *PullRequest) error {
	fullName := owner + "/" + repo

	log.Printf("DEBU: ReposQueue %s\n", ReposQueue)

	ReposQueue[fullName] = append(ReposQueue[fullName], pr)

	log.Printf("DEBU: ReposQueue %s\n", ReposQueue)

	return nil
}

// TODO: create a runner
// next will check whether or not a PR is already being processed and if not will start
// processing the next PR in the merge queue.
func Next(ctx context.Context, client *github.Client, owner, repo string) {
	fullName := owner + "/" + repo

	pr := ReposQueue[fullName][0]
	log.Printf("INFO: active PR: %s\n", pr)

	// Create a staging branch using the pull request merge commit.
	// Make sure to clean up (delete) the staging branch after each batch.
	_, err := actions.CreateStagingBranch(ctx, client, owner, repo, pr.Number)
	if err != nil {
		log.Printf("ERRO: [ CreateStagingBranch ] failed with %s\n", err)
		return
	}

	log.Println("INFO: staging branch successfully created!")
}
