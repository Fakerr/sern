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
func CreateStagingBranch(ctx context.Context, client *github.Client, owner, repo string, prNumber int) (*github.Reference, error) {

	log.Println("INFO: start [ CreateStagingBranch ]")
	defer log.Println("INFO: end [ CreateStagingBranch ]")

	stagingName := "refs/heads/" + config.StagingBranch
	log.Printf("DEBU: stagingName: %v\n", stagingName)

	sourceName := fmt.Sprintf("refs/pull/%d/merge", prNumber)
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

	// TODO:
	// Comment on the PR

	log.Println("INFO: staging branch successfully created!")

	return ref, nil
}

// Merge the pull request if the checkSuite's status is 'success'
func ProceedMerging(ctx context.Context, client *github.Client, event *github.CheckSuiteEvent, owner, repo string, activePR *queue.PullRequest) (bool, error) {

	log.Println("INFO: start [ ProceedMerging ]")
	defer log.Println("INFO: end [ ProceedMerging ]")

	activeNumber := activePR.Number
	activeSHA := activePR.MergeCommitSHA

	pr, _, err := client.PullRequests.Get(ctx, owner, repo, activeNumber)
	if err != nil {
		return false, fmt.Errorf("client.PullRequests.Get() failed with: %s\n", err)
	}

	if *pr.State != "open" {
		log.Printf("INFO: Pull request %v is no longer open\n", activeNumber)
		return false, nil
	}

	if conclusion := *event.CheckSuite.Conclusion; conclusion != "success" {
		log.Println("INFO: conclusion different from 'success', could not merge the pull request")
		// TODO:
		// Comment on the PR + update the PR's labels
		return false, nil
	}

	// TODO:
	// Comment staging branch status in the PR itself

	// Make sure the PR's SHA is the same as the active PR's SHA before merging (in case someone committed smthg during the batch)
	if activeSHA != *pr.Head.SHA {
		log.Printf("INFO: activeSHA different from the PR number %s SHA for %s/%s \n", activeNumber, owner, repo)
		// TODO:
		// Comment the problem on the PR
		return false, nil
	}

	// Merge the pullRequest
	ok := mergePullRequest(ctx, client, owner, repo, pr, activeSHA)
	if !ok {
		log.Printf("INFO: PR %v for %s/%s failed to merge!\n", activeNumber, owner, repo)
		return false, nil
	}

	log.Printf("INFO: PR %v for %s/%s merged successfully!\n", activeNumber, owner, repo)

	return true, nil
}

func mergePullRequest(ctx context.Context, client *github.Client, owner, repo string, pr *github.PullRequest, activeSHA string) bool {

	// Merge the active PR's changesets
	option := &github.PullRequestOptions{
		SHA: activeSHA,
	}

	_, _, err := client.PullRequests.Merge(ctx, owner, repo, *pr.Number, "", option)
	if err != nil {
		log.Printf("WARN: client.PullRequests.Merge() failed with %s\n", err)
		// TODO:
		// Add comment on the Pull request
		return false
	}
	return true
}
