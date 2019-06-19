package persist

import (
	"log"
)

// TEMP: simulate the db.
var enabledRepositories []*Repository

type Repository struct {
	InstallationID int64
	FullName       string
	Owner          string
	Private        bool
}

// Add a repository in the db.
func AddRepository(repo *Repository) error {
	// TEMP: simulate the db persistence. (will be added to Postgre later)
	enabledRepositories = append(enabledRepositories, repo)

	log.Printf("Enabled repos %s\n", enabledRepositories)
	return nil
}

// Remove a repository from the db.
func RemoveRepositoryByInstallationId(id int64) error {
	var tmp_repos []*Repository
	for _, elt := range enabledRepositories {
		if elt.InstallationID != id {
			tmp_repos = append(tmp_repos, elt)
		}
	}
	enabledRepositories = tmp_repos

	log.Printf("Enabled repos %s\n", enabledRepositories)
	return nil
}

// Remove a repository from the db.
func RemoveRepositoryByName(name string) error {
	var tmp_repos []*Repository
	for _, elt := range enabledRepositories {
		if elt.FullName != name {
			tmp_repos = append(tmp_repos, elt)
		}
	}
	enabledRepositories = tmp_repos

	log.Printf("Enabled repos %s\n", enabledRepositories)
	return nil
}
