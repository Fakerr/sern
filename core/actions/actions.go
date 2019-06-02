package actions

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/core/comments"
	"github.com/Fakerr/sern/core/queue"

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
	DeleteStagingBranch(ctx, client, owner, repo, stagingName)

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

	msg := "Testing with upsteam..."
	comments.AddComment(ctx, client, owner, repo, prNumber, msg)

	log.Println("INFO: staging branch successfully created!")

	return ref, nil
}

// Delete the staging branch
func DeleteStagingBranch(ctx context.Context, client *github.Client, owner, repo string, staging string) {

	log.Println("INFO: start [ DeleteStagingBranch ]")
	defer log.Println("INFO: end [ DeleteStagingBranch ]")

	_, err := client.Git.DeleteRef(ctx, owner, repo, staging)
	if err != nil {
		// If the branch doesn't exist, we can get an error
		log.Println("INFO: client.Git.DeleteRef() failed with %s but continuing...", err)
	}
}

// Merge the pull request if the checkSuite's status is 'success'
// After merge, delete the staging branch
func ProceedMerging(ctx context.Context, client *github.Client, event *github.CheckSuiteEvent, owner, repo string, activePR *queue.PullRequest) (bool, error) {

	log.Println("INFO: start [ ProceedMerging ]")
	defer log.Println("INFO: end [ ProceedMerging ]")

	activeNumber := activePR.Number
	activeSHA := activePR.HeadSHA

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
		// Comment on the PR + update the PR's labels
		msg := "Test failed"
		comments.AddComment(ctx, client, owner, repo, activeNumber, msg)
		return false, nil
	}

	// TODO:
	// Comment staging branch status in the PR itself

	// Make sure the PR's SHA is the same as the active PR's SHA before merging (in case someone committed smthg during the batch)
	if activeSHA != *pr.Head.SHA {
		log.Printf("INFO: activeSHA different from the PR number %s SHA for %s/%s \n", activeNumber, owner, repo)
		msg := "Current head" + *pr.Head.SHA + " different from accepted head " + activeSHA
		comments.AddComment(ctx, client, owner, repo, activeNumber, msg)
		return false, nil
	}

	// Merge the pullRequest
	ok := mergePullRequest(ctx, client, owner, repo, pr, activeSHA)
	if !ok {
		log.Printf("INFO: PR %v for %s/%s failed to merge!\n", activeNumber, owner, repo)
		return false, nil
	}

	log.Printf("INFO: PR %v for %s/%s merged successfully!\n", activeNumber, owner, repo)
	log.Printf("INFO: Deleting staging branch for %s/%s\n", owner, repo)

	stagingName := "refs/heads/" + config.StagingBranch
	DeleteStagingBranch(ctx, client, owner, repo, stagingName)

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
		// Items that failed to merge will be ignored and removed from the queue (for now)
		return false
	}
	return true
}

// Check whether or not a PR is still meargeable
func CheckMergeability(ctx context.Context, client *github.Client, owner, repo string, num int, pr *github.PullRequest) bool {

	log.Printf("INFO: start checking mergeability for %s/%s number: %v\n", owner, repo, num)

	mergeable := pr.Mergeable

	if mergeable == nil {
		// If the value is nil, then GitHub has started a background job to compute the mergeability and it's not complete yet.
		// Sleep 5 seconds and try again
		log.Printf("INFO: merageability info not yet ready for %s/%s number: %v, sleeping for 5 seconds...\n", owner, repo, num)
		time.Sleep(5 * time.Second)

		pr, _, err := client.PullRequests.Get(ctx, owner, repo, num)
		if err != nil {
			log.Printf("ERRO: client.PullRequests.Get() failed for %s/%s number: %v with: %s\n", owner, repo, num, err)
			return false
		}

		if pr.Mergeable == nil {
			log.Printf("INFO: Cannot get merageability info for %s/%s number: %v but treating it as if mergeable\n", owner, repo, num)
			return true
		}

		return *pr.Mergeable
	}

	return *pr.Mergeable
}
