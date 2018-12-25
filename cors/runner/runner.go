package runner

import (
	"context"
	"log"

	"github.com/Fakerr/sern/cors/actions"
	"github.com/Fakerr/sern/cors/queue"

	"github.com/google/go-github/github"
)

// A map with a all enabled repositories and their current runner. Should be presisted (redis).
var ReposRunner map[string]*Runner

func SetRunners() {
	ReposRunner = make(map[string]*Runner)
}

// If it exists, return the owner/repo's runner, otherwise create a new one.
func GetRunner(owner, repo string) *Runner {
	name := owner + "/" + repo

	if ReposRunner[name] != nil {
		return ReposRunner[name]
	}

	runner := &Runner{
		Owner:  owner,
		Repo:   repo,
		Status: "inactive", // running or inactive
		Locked: false,
		Active: nil,
		Queue:  &queue.QueueMerge{},
	}
	initQueue := make([]*queue.PullRequest, 0)
	runner.Queue.Init(initQueue)

	return runner
}

type Runner struct {
	Owner  string
	Repo   string
	Status string // running or inactive
	Locked bool
	Active *queue.PullRequest
	Queue  *queue.QueueMerge
}

func (r *Runner) SetStatus(status string) {
	r.Status = status
}

func (r *Runner) RemoveActive() {
	r.Active = nil
	r.SetStatus("inactive")
	r.Queue.RemoveFirst()
}

// Next will check whether or not a PR is already being processed and if not will start
// processing the next PR in the merge queue.
func (r *Runner) Next(ctx context.Context, client *github.Client) {

	log.Printf("INFO: start [ runner.Next ] for %s/%s \n", r.Owner, r.Repo)
	defer log.Printf("INFO: end [ runner.Next ] for %s/%s \n", r.Owner, r.Repo)

	// If already running, don't do anything
	if r.Status == "running" || r.Locked {
		return
	}

	r.Locked = true
	nextItem := r.getNextItem()

	// If no item left in the queue.
	if nextItem == nil {
		r.Locked = false
		return
	}

	r.Active = nextItem
	r.SetStatus("running")

	// Create a staging branch using the pull request merge commit.
	_, err := actions.CreateStagingBranch(ctx, client, r.Owner, r.Repo, r.Active.Number)
	if err != nil {
		log.Printf("ERRO: [ CreateStagingBranch ] failed with %s\n", err)
		log.Println("INFO: trying another item...")

		r.Locked = false
		r.RemoveActive()
		r.Next(ctx, client)
		return
	}

	r.Locked = false
	return
}

func (r *Runner) getNextItem() *queue.PullRequest {

	next := r.Queue.GetFirst()

	// If not item left in the queue, return nil
	if next == nil {
		return nil
	}

	// Before returning the next Item, make sure the accepted PR is still accurate (PR still open, same SHA,...)
	// otherwise, ignore and take the item that comes after

	return next
}
