package persist

import (
	"log"
)

// TEMP: simulate the db.
var enabledRepositories []*Repository

type Repository struct {
	Name  string
	Owner string
	// The first time the user enable a repo, it should be presisted with its config.
	// Later the user can disable the repo but its config is still saved.
	Enabled bool
}

// Add a repository in the db.
func AddRepository(repo *Repository) error {
	// TEMP: simulate the db persistence. (will be added to Postgre later)
	enabledRepositories := append(enabledRepositories, repo)
	log.Printf("Enabled repos %s\n", enabledRepositories)
	return nil
}

// Remove a repository from the db.
func RemoveRepository(repoID string) {

}
