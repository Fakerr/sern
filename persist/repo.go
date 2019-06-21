package persist

import (
	"log"
)

// TEMP: simulate the db.
var EnabledRepositories []*Repository

type Repository struct {
	InstallationID int64
	FullName       string
	Owner          string
	Private        bool
	MergeType      string // Can be either empty, "squash" or "rebase"
}

func GetRepositoryByName(name string) *Repository {
	for _, elt := range EnabledRepositories {
		if elt.FullName == name {
			return elt
		}
	}
	return nil
}

func GetRepositoriesByOwner(owner string) (repos []*Repository) {
	for _, elt := range EnabledRepositories {
		if elt.Owner == owner {
			repos = append(repos, elt)
		}
	}
	return
}

// Return true if the repository is a private repository
func IsPrivate(name string) bool {
	repo := GetRepositoryByName(name)
	if repo != nil {
		return repo.Private
	}
	return false
}

// Add a repository in the db.
func AddRepository(repo *Repository) error {
	// TEMP: simulate the db persistence. (will be added to Postgre later)
	EnabledRepositories = append(EnabledRepositories, repo)

	log.Printf("Enabled repos %s\n", EnabledRepositories)
	return nil
}

// Remove a repository from the db.
func RemoveRepositoryByInstallationId(id int64) error {
	var tmp_repos []*Repository
	for _, elt := range EnabledRepositories {
		if elt.InstallationID != id {
			tmp_repos = append(tmp_repos, elt)
		}
	}
	EnabledRepositories = tmp_repos

	log.Printf("Enabled repos %s\n", EnabledRepositories)
	return nil
}

// Remove a repository from the db.
func RemoveRepositoryByName(name string) error {
	var tmp_repos []*Repository
	for _, elt := range EnabledRepositories {
		if elt.FullName != name {
			tmp_repos = append(tmp_repos, elt)
		}
	}
	EnabledRepositories = tmp_repos

	log.Printf("Enabled repos %s\n", EnabledRepositories)
	return nil
}
