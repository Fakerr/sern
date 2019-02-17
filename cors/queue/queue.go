package queue

import (
	"context"
	"log"

	"github.com/Fakerr/sern/cors/comments"

	"github.com/google/go-github/github"
)

type PullRequest struct {
	Number         int
	HeadSHA        string
	MergeCommitSHA string
}

type QueueMerge struct {
	QueueItems []*PullRequest
}

// Init the merge queue
func (q *QueueMerge) Init(queue []*PullRequest) {
	q.QueueItems = queue
}

// Get first item of the merge queue.
func (q *QueueMerge) GetFirst() *PullRequest {
	if len(q.QueueItems) == 0 {
		return nil
	}
	return q.QueueItems[0]
}

// Delete first item from the merge queue.
func (q *QueueMerge) RemoveFirst() {
	q.QueueItems = q.QueueItems[1:]
}

// Add Pull Request to the merge queue.
func (q *QueueMerge) Add(ctx context.Context, client *github.Client, owner, repo string, pr *PullRequest) {
	//make sure the PR is not already in the queue.
	for _, item := range q.QueueItems {
		if item.Number == pr.Number {
			log.Printf("WARN: Pull Request %v already queued", pr)
			// msg := "PR number " + pr.Number + " already queued."
			msg := "msg 1"
			comments.AddComment(ctx, client, owner, repo, pr.Number, msg)
			return
		}
	}
	q.QueueItems = append(q.QueueItems, pr)
	// msg := "Pull request added in the merge queue."
	msg := "msg 2"
	comments.AddComment(ctx, client, owner, repo, pr.Number, msg)
	return
}

// Remove Pull Request from the merge queue.
func (q *QueueMerge) Remove(pr *PullRequest) {
}
