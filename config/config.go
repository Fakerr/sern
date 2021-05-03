package config

import (
	"os"
	"strconv"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// The app must be registered at https://github.com/settings/applications
// Set callback to http://127.0.0.1:8089/github_oauth_cb
// Set ClientId and ClientSecret
// Set webhooksecret

var IntegrationID, _ = strconv.Atoi(os.Getenv("INTEGRATION_ID"))
var PrivateKey = os.Getenv("PRIVATE_KEY")
var WebHookSecret = os.Getenv("WEBHOOK_SECRET") // Set within the Github App configuration on github
var SessionSecretKey = os.Getenv("SESSION_SECRET_KEY")
var oauthClientID = os.Getenv("OAUTH_CLIENT_ID")
var oauthClientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
var DatabaseURL = os.Getenv("DATABASE_URL")
var RedisURL = os.Getenv("REDIS_URL")

const STAGING_BRANCH = "test-staging"
const CMD_MERGE = "sern-merge"

// Github oauth configuration
var OauthConf = &oauth2.Config{
	ClientID:     oauthClientID,
	ClientSecret: oauthClientSecret,
	Endpoint:     githuboauth.Endpoint,
}
