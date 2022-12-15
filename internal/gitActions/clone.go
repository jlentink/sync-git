package gitActions

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/subchen/go-log"
	"os"
	"os/exec"
	"strings"
	"sync-git/internal/cnf"
)

var _repository *git.Repository

func HasClone() bool {
	source := cnf.GetSource()
	if _, err := os.Stat(source.BaseName()); err == nil {
		log.Debugf("Found %s. Skipping...", getBaseName())
		return true
	}
	return false
}

func Clone() {
	source := cnf.GetSource()
	if HasClone() {
		return
	}

	log.Infof("Cloning %s", source.BaseName())
	repo, err := git.PlainClone("./"+source.BaseName(), false, &git.CloneOptions{
		URL:      source.CloneUrl(),
		Progress: os.Stdout,
	})
	_repository = repo
	if err != nil {
		log.Fatalf("Error cloning %s", err)
	}

}

func AddRemotes() {
	repo := Open()
	remotes, err := repo.Remotes()
	if err != nil {
		log.Fatalf("Error getting remotes %s", err)
	}

	for _, destination := range cnf.GetDestinations() {
		found := false
		for _, remote := range remotes {
			if remote.Config().Name == destination.RemoteName() {
				found = true
			}
		}
		if !found {
			log.Debugf("Adding remote %s", destination.RemoteName())
			_, err := _repository.CreateRemote(&config.RemoteConfig{
				Name: destination.RemoteName(),
				URLs: []string{destination.CloneUrl()},
			})
			if err != nil {
				log.Fatalf("Error adding remote %s", err)
			}
		}
	}

}

func Open() *git.Repository {
	if _repository == nil {
		repo, err := git.PlainOpen(getBaseName())
		if err != nil {
			log.Fatalf("Error opening %s", err)
		}
		_repository = repo

	}
	return _repository
}

func Fetch() {
	err := Open().Fetch(&git.FetchOptions{})
	if err != nil {
		if err.Error() != "already up-to-date" {
			log.Fatalf("Error fetching %s", err)
		}
	}
}

func execute(args ...string) ([]string, error) {

	cmd := exec.Command(cnf.GetString("git.path")+"/git", args...)
	cmd.Dir = cnf.GetSource().BaseName()
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}
