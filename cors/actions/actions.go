package actions

import (
	"context"
	"fmt"
	"log"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/cors/queue"

	"github.com/google/go-github/github"
)

// Create a staging branch to test branch with upstream
func CreateStagingBranch(ctx context.Context, client *github.Client, owner, repo string, prID int) (*github.Reference, error) {

	log.Println("INFO: start [ CreateStagingBranch ]")
	defer log.Println("INFO: end [ CreateStagingBranch ]")

	stagingName := "refs/heads/" + config.StagingBranch
	log.Printf("DEBU: stagingName: %v\n", stagingName)

	sourceName := fmt.Sprintf("refs/pull/%d/merge", prID)
	log.Printf("DEBU: sourceName: %v\n", sourceName)

	// Clean up staging branch before processing
	log.Println("INFO: clean up staging branch")
	if _, err := client.Git.DeleteRef(ctx, owner, repo, stagingName); err != nil {
		log.Println("INFO: client.Git.DeleteRef() failed with %s but continuing...", err)
	}

	sourceRef, _, err := client.Git.GetRef(ctx, owner, repo, sourceName)
	if err != nil {
		return nil, fmt.Errorf("client.Git.GetRef() failed for %v with: %s\n", sourceName, err)
	}

	stagingRef := github.Reference{
		Ref:    &stagingName,
		URL:    nil,
		Object: sourceRef.Object,
	}

	ref, _, err := client.Git.CreateRef(ctx, owner, repo, &stagingRef)
	if err != nil {
		return nil, fmt.Errorf("client.Git.CreateRef() failed for %v with: %s\n", stagingName, err)
	}

	return ref, nil
}

// Merge the pull request if the checkSuite's status is 'success'
func ProceedMerging(ctx context, client *github.Client, event *github.CheckSuiteEvent, owner, repo string, activePR *queue.PullRequest) error {

	log.Println("INFO: start [ ProceedMerging ]")
	defer log.Println("INFO: end [ ProceedMerging ]")

	prID := activePR.Id

	prInfo, _, err := client.PullRequests.Get(ctx, owner, repo, prId)
	if err != nil {
		return fmt.Errorf("client.PullRequests.Get() failed for with: %s\n", err)
	}

	if *prInfo.State != "open" {
		log.Printf("INFOR: Pull request %v is no longer open\n", prID)
		return nil
	}

	if conclusion := *event.CheckSuite.Conclusion; conclusion != "success" {
		log.Println("INFO: conclusion different from 'success', could not merge the pull request")
		// Comment on the PR + update the PR's labels
		return nil
	}

	// Comment staging branch status.

	// Merge the pullRequest

	log.Println("INFO: PR %v for %s/%s merged successfully!\n", prID, owner, repo)
}
