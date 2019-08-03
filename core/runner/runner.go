package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Fakerr/sern/core/actions"
	"github.com/Fakerr/sern/core/comments"
	"github.com/Fakerr/sern/core/queue"
	"github.com/Fakerr/sern/persist"

	"github.com/gomodule/redigo/redis"
	"github.com/google/go-github/github"
)

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

// update runner in database
func (r *Runner) Update() {
	err := setRunner(r)
	if err != nil {
		log.Printf("ERRO: [ update ] failed with %s\n", err)
	}
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

	// updating the runner in db after the lock to prevent any concurrent transaction
	r.Update()

	nextItem := r.getNextItem(ctx, client)

	// If no item left in the queue, destroy the runner
	if nextItem == nil {
		name := r.Owner + "/" + r.Repo
		log.Printf("INFO: no items left in the queue for %s, deleting the runner...\n", name)
		err := deleteRunner(name)
		if err != nil {
			log.Printf("ERRO: [ deleteRunner ] for %s . Failed with %s\n", name, err)
		}
		return
	}

	r.Active = nextItem
	r.SetStatus("running")

	// Create a staging branch using the pull request merge commit.
	ref, err := actions.CreateStagingBranch(ctx, client, r.Owner, r.Repo, r.Active.Number)
	if err != nil {
		log.Printf("ERRO: [ CreateStagingBranch ] failed with %s\n", err)
		log.Println("INFO: trying another item...")

		r.RemoveActive()
		r.Locked = false
		r.Update()
		r.Next(ctx, client)
		return
	}

	// Update the activePR merge commit's SHA (the merge commmit sha change each time a new change is added into upstream)
	r.Active.MergeCommitSHA = *ref.Object.SHA
	r.Locked = false
	r.Update()

	return
}

func (r *Runner) getNextItem(ctx context.Context, client *github.Client) *queue.PullRequest {

	log.Printf("INFO: start [ runner.getNextItem ] for %s/%s \n", r.Owner, r.Repo)
	defer log.Printf("INFO: end [ runner.getNextItem ] for %s/%s \n", r.Owner, r.Repo)

	for {
		next := r.Queue.GetFirst()

		// If not item left in the queue, return nil
		if next == nil {
			return nil
		}

		num := next.Number

		pr, _, err := client.PullRequests.Get(ctx, r.Owner, r.Repo, num)
		if err != nil {
			log.Printf("ERRO: client.PullRequests.Get() failed for %s/%s number: %v with: %s\n", r.Owner, r.Repo, num, err)
			r.Queue.RemoveFirst()
			continue
		}

		if state := *pr.State; state != "open" {
			log.Printf("INFO: Pull request %v is no longer open for %s/%s\n", num, r.Owner, r.Repo)
			r.Queue.RemoveFirst()
			continue
		}

		if next.HeadSHA != *pr.Head.SHA {
			log.Printf("INFO: next's SHA different from the PR number: %s SHA for %s/%s \n", num, r.Owner, r.Repo)
			r.Queue.RemoveFirst()
			msg := "Current head" + *pr.Head.SHA + " different from accepted head " + next.HeadSHA
			comments.AddComment(ctx, client, r.Owner, r.Repo, num, msg)
			continue
		}

		mergeable := actions.CheckMergeability(ctx, client, r.Owner, r.Repo, num, pr)
		if !mergeable {
			log.Printf("INFO: PR for %s/%s number: %v is not mergeable\n", r.Owner, r.Repo, num)
			r.Queue.RemoveFirst()
			msg := "Merge conflict! please resolve"
			comments.AddComment(ctx, client, r.Owner, r.Repo, num, msg)
			continue
		}

		// TODO:
		// Check if all required labels are set correctly

		return next
	}
}

// If it exists, return the owner/repo's runner, otherwise create a new one.
// If no items left, the runner should be destroyed
func GetRunner(owner, repo string) *Runner {
	name := owner + "/" + repo

	runner, err := getRunnerFromdb(name)
	if err != nil {
		log.Printf("ERRO: [ getRunnerFromdb ] failed with %s\n", err)
		return nil
	}
	if runner != nil {
		return runner
	}

	runner = &Runner{
		Owner:  owner,
		Repo:   repo,
		Status: "inactive", // running or inactive
		Locked: false,
		Active: nil,
		Queue:  &queue.QueueMerge{},
	}

	initQueue := make([]*queue.PullRequest, 0)
	runner.Queue.Init(initQueue)

	// Persist the new runner in database
	err = setRunner(runner)
	if err != nil {
		log.Printf("ERRO: [ setRunner ] failed with %s\n", err)
		return nil
	}

	log.Printf("INFO: A new runner for %s has been created!\n", name)

	return runner
}

// Get the runner without creating a new one if it doesn't exist
func GetSoftRunner(owner, repo string) *Runner {
	name := owner + "/" + repo

	runner, err := getRunnerFromdb(name)
	if err != nil {
		log.Printf("ERRO: [ getRunnerFromdb ] failed with %s\n", err)
		return nil
	}
	if runner != nil {
		return runner
	}

	return nil
}

// Persist runner in the redis db
func setRunner(runner *Runner) error {

	conn := persist.Pool.Get()
	defer conn.Close()

	var key string = runner.Owner + "/" + runner.Repo

	// serialize object to JSON
	json, err := json.Marshal(runner)
	if err != nil {
		return err
	}

	// SET object
	_, err = conn.Do("SET", key, json)
	if err != nil {
		return err
	}

	return nil
}

// Get runner from redis
func getRunnerFromdb(key string) (*Runner, error) {

	conn := persist.Pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))

	if err == redis.ErrNil {
		fmt.Printf("Runner %s does not exist", key)
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	runner := Runner{}
	err = json.Unmarshal([]byte(s), &runner)

	fmt.Printf("%+v\n", runner)

	return &runner, nil
}

// delete runner from redis
func deleteRunner(key string) error {

	conn := persist.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}
