package hooks

import (
	"context"
	"fmt"
	"log"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/core/actions"
	"github.com/Fakerr/sern/core/client"
	"github.com/Fakerr/sern/core/queue"
	"github.com/Fakerr/sern/core/runner"
	"github.com/Fakerr/sern/persist"

	"github.com/google/go-github/github"
)

// Handle IssueComment event
func ProcessIssueCommentEvent(ctx context.Context, event *github.IssueCommentEvent) error {

	log.Println("INFO: start [ ProcessIssueCommentEvent ]")
	defer log.Println("INFO: end [ ProcessIssueCommentEvent ]")

	// Return an error if the action is different from "created"
	if action := event.Action; (action == nil) || (*action != "created") {
		return fmt.Errorf("Accept only `action === \"created\"` event")
	}

	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name

	log.Printf("INFO: processing %s/%s \n", owner, repo)

	// Check whether or not the Issue Comment was made on a Pull Request.
	// If not, return as nothing to do.
	if event.Issue.IsPullRequest() == false {
		return nil
	}

	// if the Issue Comment is not a valid command, return
	cmd, ok := parseComment(*event.Comment.Body)
	if !ok {
		log.Println("INFO: aborting: not a sern command.")
		return nil
	}

	// For now, make the commands available only for the repo's owner.
	if owner != *event.Comment.User.Login {
		return fmt.Errorf("Only the repo's owner is able to run this command.")
	}

	// Create an installation client.
	client := client.GetInstallationClient(int(*event.Installation.ID))

	// Take appropriate action based on the parsed command.
	switch cmd {
	case config.CMD_MERGE:
		pr, ok, err := createPullRequest(ctx, client, owner, repo, event)

		if err != nil {
			return fmt.Errorf("[ createPullRequest ] failed with %s\n", err)
		}

		if !ok {
			log.Println("INFO: PR's upstream branch different from master.")
			return nil
		}

		// Get the current runner or create a new one if it doesn't exist
		runner := runner.GetRunner(owner, repo)
		if runner != nil {
			runner.Queue.Add(ctx, client, owner, repo, pr)
			runner.Update()
			runner.Next(ctx, client)
		}
	default:
		log.Printf("INFO: No handler for command: %v.", cmd)
	}

	return nil
}

// Handle CheckSuite event
func ProcessCheckSuiteEvent(ctx context.Context, event *github.CheckSuiteEvent) error {

	log.Println("INFO: start [ ProcessCheckSuiteEvent ]")
	defer log.Println("INFO: end [ ProcessCheckSuiteEvent ]")

	fullName := *event.Repo.FullName
	log.Printf("INFO: processing repository: %s\n", fullName)

	// Make sure the CheckSuite event is about the staging branch.
	log.Printf("DEBU: checkSuite headBranch: %s\n", *event.CheckSuite.HeadBranch)
	if *event.CheckSuite.HeadBranch != config.STAGING_BRANCH {
		log.Println("INFO: CheckSuite's headBranch different from StagingBranch. Aborting!")
		return nil
	}

	log.Printf("DEBU: CheckSuite status: %s\n", *event.CheckSuite.Status)
	if status := *event.CheckSuite.Status; status != "completed" {
		log.Println("INFO: status different from 'completed' is not handled")
		return nil
	}

	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name

	// Get the runner instance
	runner := runner.GetRunner(owner, repo)
	if runner == nil {
		return nil
	}
	activePR := runner.Active

	// Make sure the event's commit hash is the same as the active PR's merge commit hash.
	if activePR.MergeCommitSHA != *event.CheckSuite.HeadSHA {
		log.Println("INFO: event's commit hash different from the active PR's merge commit hash")
		return nil
	}

	// Create an installation client.
	client := client.GetInstallationClient(int(*event.Installation.ID))

	_, err := actions.ProceedMerging(ctx, client, event, owner, repo, activePR)
	if err != nil {
		return fmt.Errorf("[ actions.ProceedMerging ] failed with %s\n", err)
	}

	// Regardless the previous item succeed to merge or not, proceed to the next item
	runner.RemoveActive()
	runner.Update()
	runner.Next(ctx, client)

	return nil
}

// Handler for Installation event
func ProcessInstallationEvent(ctx context.Context, event *github.InstallationEvent) error {

	log.Println("INFO: start [ ProcessInstallationEvent ]")
	defer log.Println("INFO: end [ ProcessInstallationEvent ]")

	if *event.Action == "created" {
		for _, item := range event.Repositories {
			// Create the repository and presist it in the db.
			repository := &persist.Repository{
				InstallationID: *event.Installation.ID,
				FullName:       *item.FullName,
				Owner:          *event.Sender.Login,
				Private:        *item.Private,
			}

			// Enable repository (Persist in the db).
			err := persist.AddRepository(repository)

			if err != nil {
				return fmt.Errorf("[ persist.AddRepository ] failed with %s\n", err)
			}
		}
	}

	if *event.Action == "deleted" {
		err := persist.RemoveRepositoryByInstallationId(*event.Installation.ID)
		if err != nil {
			return fmt.Errorf("[ persist.RemoveRepositoryByInstallationId ] failed with %s\n", err)
		}
	}

	return nil
}

// Handler for InstallationRepositories event
func ProcessInstallationRepositoriesEvent(ctx context.Context, event *github.InstallationRepositoriesEvent) error {

	log.Println("INFO: start [ ProcessInstallationRepositoriesEvent ]")
	defer log.Println("INFO: end [ ProcessInstallationRepositoriesEvent ]")

	if *event.Action == "added" {
		for _, item := range event.RepositoriesAdded {
			// Create the repository and presist it in the db.
			repository := &persist.Repository{
				InstallationID: *event.Installation.ID,
				FullName:       *item.FullName,
				Owner:          *event.Sender.Login,
				Private:        *item.Private,
			}

			// Enable repository (Persist in the db).
			err := persist.AddRepository(repository)

			if err != nil {
				return fmt.Errorf("[ persist.AddRepository ] failed with %s\n", err)
			}
		}
	}

	if *event.Action == "removed" {
		for _, item := range event.RepositoriesRemoved {
			err := persist.RemoveRepositoryByName(*item.FullName)
			if err != nil {
				return fmt.Errorf("[ persist.RemoveRepositoryByName ] failed with %s\n", err)
			}
		}
	}

	return nil
}

// Parses the comment body and return an action
func parseComment(body string) (string, bool) {
	switch body {
	case config.CMD_MERGE:
		return config.CMD_MERGE, true
	default:
		return "", false
	}
}

// Create a new pull request ready to be queued.
// If the pull request's upsteam branch is not master, ok will be set to false and the PR won't be queued.
func createPullRequest(ctx context.Context, client *github.Client, owner, repo string, event *github.IssueCommentEvent) (*queue.PullRequest, bool, error) {

	log.Println("INFO: start [ createPullRequest ]")
	defer log.Println("INFO: end [ createPullRequest ]")

	number := *event.Issue.Number
	fullName := owner + "/" + repo

	// Fetch PullRequest
	pull, _, err := client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, false, fmt.Errorf("client.PullRequests.Get() failed for %s with: %s\n", fullName, err)
	}

	if *pull.Base.Ref != "master" {
		return nil, false, nil
	}

	pr := &queue.PullRequest{
		Number:         number,
		HeadSHA:        *pull.Head.SHA,
		MergeCommitSHA: *pull.MergeCommitSHA,
	}
	return pr, true, nil
}
