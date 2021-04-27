package persist

import ()

type Repository struct {
	InstallationID int64
	FullName       string
	Owner          string
	Private        bool
}

func GetRepositoriesByOwner(owner string) ([]*Repository, error) {
	rows, err := Conn.Query("SELECT installation, fullname, owner, private FROM repositories WHERE owner = $1;", owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repositories := make([]*Repository, 0)
	for rows.Next() {
		repo := Repository{}
		err := rows.Scan(&repo.InstallationID, &repo.FullName, &repo.Owner, &repo.Private)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, &repo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return repositories, nil
}

func GetRepositoryByName(name string) (*Repository, error) {
	rows, err := Conn.Query("SELECT installation, fullname, owner, private FROM repositories WHERE fullname = $1;", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	repo := Repository{}
	for rows.Next() {
		err := rows.Scan(&repo.InstallationID, &repo.FullName, &repo.Owner, &repo.Private)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &repo, nil
}

// Return true if repository is installed
func RepositoryExists(fullname string) (bool, error) {
	exist := false
	rows, err := Conn.Query("SELECT * FROM repositories WHERE fullname = $1;", fullname)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		//some values exist
		exist = true
		break
	}
	return exist, nil
}

// Return true if the repository is a private repository
func IsPrivate(name string) (bool, error) {
	repo, err := GetRepositoryByName(name)
	if err != nil {
		return false, err
	}

	if repo != nil {
		return repo.Private, nil
	}
	return false, nil
}

// Add a repository in the db.
func AddRepository(repo *Repository) error {
	query := "INSERT INTO repositories (installation, fullname, owner, private) VALUES ($1, $2, $3, $4);"
	_, err := Conn.Exec(query, repo.InstallationID, repo.FullName, repo.Owner, repo.Private)
	if err != nil {
		return err
	}
	return nil
}

// Remove a repository from the db.
func RemoveRepositoryByInstallationId(id int64) error {
	_, err := Conn.Exec("DELETE FROM repositories WHERE installation=$1;", id)
	if err != nil {
		return err
	}
	return nil
}

// Remove a repository from the db.
func RemoveRepositoryByName(name string) error {
	_, err := Conn.Exec("DELETE FROM repositories WHERE fullname=$1;", name)
	if err != nil {
		return err
	}
	return nil
}

func GetInstallationIDByName(name string) (*int, error) {
	var installation int
	row := Conn.QueryRow("SELECT installation FROM repositories WHERE fullname = $1;", name)
	err := row.Scan(&installation)
	if err != nil {
		return nil, err
	}
	return &installation, nil
}
