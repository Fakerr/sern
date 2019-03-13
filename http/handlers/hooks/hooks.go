package hooks

import (
	"log"
	"net/http"
	"reflect"

	"github.com/Fakerr/sern/config"
	"github.com/Fakerr/sern/cors/hooks"

	"github.com/google/go-github/github"
)

// Endpoint that will receive the webhook POST requests
func WebhookCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := github.ValidatePayload(r, []byte(config.WebHookSecret))
	if err != nil {
		log.Printf("github.ValidatePayload() failed with '%s'\n", err)
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("github.ParseWebHook() failed with '%s'\n", err)
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
		return
	}

	switch event := event.(type) {
	case *github.IssueCommentEvent:
		err := hooks.ProcessIssueCommentEvent(ctx, event)

		if err != nil {
			log.Printf("ERRO: [ ProcessIssueCommentEvent ] failed with: %s\n", err)
		}

		w.WriteHeader(http.StatusOK)
		return
	case *github.CheckSuiteEvent:
		err := hooks.ProcessCheckSuiteEvent(ctx, event)

		if err != nil {
			log.Printf("ERRO: [ ProcessCheckSuiteEvent ] failed with: %s\n", err)
		}

		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusOK)
		log.Printf("WARN: Unsupported type events %v\n", reflect.TypeOf(event))
		//io.WriteString(rw, "This event type is not supported: "+github.WebHookType(req))
		return
	}
}
