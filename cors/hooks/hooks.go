package hooks

import (
	"context"
	"fmt"
	"log"

	"github.com/Fakerr/sern/cors/queue"

	"github.com/google/go-github/github"
)

// Handle IssueComment event
func ProcessIssueCommentEvent(ctx context.Context, event *github.IssueCommentEvent) error {

	// Return an error if the action is different from "created"
	if action := event.Action; (action == nil) || (*action != "created") {
		return fmt.Errorf("Accept only `action === \"created\"` event")
	}

	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name
	fullRepo := owner + "/" + repo

	log.Printf("ProcessIssueCommentEvent %s\n", fullRepo)

	body := *event.Comment.Body

	// Check whether or not the Issue Comment was made on a Pull Request.
	// If not, return as nothing to do.
	if event.Issue.IsPullRequest() == false {
		return nil
	}

	// if the Issue Comment is not valid, return
	cmd, ok := parseComment(body)
	if !ok {
		return nil
	}

	if cmd == "test" {
		pr, err := createPullRequest(event)
		if err != nil {
			return fmt.Errorf("createPullRequest error %s", err)
		}

		err = queue.AddToQueue(fullRepo, pr)
		if err != nil {
			return fmt.Errorf("queue.AddToQueue error %s", err)
		}
	}

	return nil
}

// Parses the comment body and return an action
func parseComment(body string) (string, bool) {
	if body == "test" {
		return "test", true
	}
	return "", false
}

// Create a new pull request ready to be queued.
func createPullRequest(event *github.IssueCommentEvent) (*queue.PullRequest, error) {
	pr := &queue.PullRequest{
		Id:     *event.Issue.Number,
		Status: "pending",
	}
	return pr, nil
}
