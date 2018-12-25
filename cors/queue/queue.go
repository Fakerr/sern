package queue

import (
	"log"
)

type PullRequest struct {
	Number         int
	Status         string
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
func (q *QueueMerge) Add(pr *PullRequest) {
	//make sure the PR is not already in the queue.
	for _, item := range q.QueueItems {
		if item.Number == pr.Number {
			log.Printf("WARN: Pull Request %v already queued", pr)
			// TODO:
			// Add comment on the PR (already queued)
			return
		}
	}
	q.QueueItems = append(q.QueueItems, pr)
	return
}

// Remove Pull Request from the merge queue.
func (q *QueueMerge) Remove(pr *PullRequest) {
}
