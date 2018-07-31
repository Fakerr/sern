package client

import (
	"context"
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

// Get a new github client using an access token
func FromToken(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.TODO(), ts)

	return github.NewClient(tc)
}
