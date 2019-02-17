package comments

import (
	"context"
	"log"

	"github.com/google/go-github/github"
)

// Add comment on Pull request
func AddComment(ctx context.Context, client *github.Client, owner, repo string, num int, msg string) {

	log.Printf("INFO: Adding comment for %s/%s number: %v\n", owner, repo, num)

	_, _, err := client.Issues.CreateComment(ctx, owner, repo, num, &github.IssueComment{
		Body: &msg,
	})

	if err != nil {
		log.Printf("ERRO: client.Issues.CreateComment() failed for %s/%s#%v with: %s\n", owner, repo, num, err)
	}
}
