package client

import (
	"context"
	"golang.org/x/oauth2"
	"log"
	"net/http"

	"github.com/Fakerr/sern/config"

	"github.com/bradleyfalzon/ghinstallation"
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

// Return an installation client.
func GetInstallationClient(installationID int) *github.Client {
	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport

	// Wrap the shared transport for use with the integration ID authenticating with installation ID.
	itr, err := ghinstallation.NewKeyFromFile(tr, config.IntegrationID, installationID, config.PrivateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Use installation transport with github.com/google/go-github
	return github.NewClient(&http.Client{Transport: itr})
}
