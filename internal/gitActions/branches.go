package gitActions

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/subchen/go-log"
	"os"
	"strings"
	"sync-git/internal/cnf"
)

func getBranches(remote bool) []string {
	options := []string{"branch", "--format='%(refname:short)'"}
	if remote {
		options = append(options, "-r")
	}
	output, err := execute(options...)
	branches := make([]string, 0)
	if err != nil {
		log.Errorf("Error listing branches %s (%t)", err, remote)
		os.Exit(1)
	}
	for _, branch := range output {
		branch = strings.TrimSpace(branch)
		branch = strings.Replace(branch, "'", "", -1)
		if len(branch) > 0 {
			branches = append(branches, branch)
		}
	}
	return branches

}

func ListRemoteBranches(remote string) []string {
	branches := getBranches(true)
	if len(remote) == 0 {
		return branches
	}
	filtered := make([]string, 0)
	for _, branch := range branches {
		if strings.HasPrefix(branch, remote) {
			filtered = append(filtered, branch)
		}
	}
	return filtered

}
func ListLocalBranches() []string {
	return getBranches(false)
}

func BranchExists(branch string) bool {
	branches := ListLocalBranches()
	for _, b := range branches {
		if b == branch {
			return true
		}
	}
	return false
}

func Checkout(branch string) {
	log.Debugf("Checking out %s", branch)
	if BranchExists(branch) {
		execute("checkout", "--track", branch) //nolint:errcheck
	} else {
		localBranch := strings.Replace(branch, "origin/", "", -1)
		execute("checkout", localBranch) //nolint:errcheck
	}

}

func Pull(branch string) {
	localBranch := strings.Replace(branch, "origin/", "", -1)
	log.Debugf("Pulling %s", localBranch)
	execute("pull", "origin", localBranch) //nolint:errcheck
}

func Push(destination *cnf.GitDestination, branch string) {
	repo := Open()
	log.Debugf("Pushing %s to %s", branch, destination.RemoteName())
	err := repo.Push(&git.PushOptions{
		RemoteURL: destination.CloneUrl(),
		Force:     true,
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Debugf("Branch %s is up to date", branch)
		} else {
			log.Errorf("Error pushing %s", err)
			os.Exit(1)
		}
	}
}

func DeleteArtifactBranches(destination *cnf.GitDestination) {
	execute("fetch", "--all", "-p") //nolint:errcheck
	sourceBranches := StripRemote(ListRemoteBranches("origin"), "origin/")
	remoteBranches := StripRemote(ListRemoteBranches(destination.RemoteName()), destination.RemoteName())
	for _, branch := range remoteBranches {
		if !contains(sourceBranches, branch) {
			log.Infof("Deleting unused branch %s in remote %s", branch, destination.RemoteName())
			DeleteBranch(branch, destination.RemoteName())
		}
	}
}

func StripRemote(branches []string, remote string) []string {
	for i, branch := range branches {
		branches[i] = branch[strings.Index(branch, "/")+1:]
	}
	return branches
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func DeleteBranch(branch, remote string) {
	log.Debugf("Deleting branch %s", branch)
	execute("branch", "-D", "-r", fmt.Sprintf("%s/%s", remote, branch)) //nolint:errcheck
	execute("push", remote, ":"+branch)                                 //nolint:errcheck
	execute("branch", "-D", branch)                                     //nolint:errcheck
}

func CompareRemoteChanged(branch string) bool {
	localBranch := strings.Replace(branch, "origin/", "", -1)
	if !BranchExists(localBranch) {
		return true
	}
	remoteHash := BranchHash(branch)
	localHash := BranchHash(localBranch)
	if remoteHash == localHash {
		log.Infof("Branch %s (%s) is up to date", branch, remoteHash)
		return false
	}
	return true
}

func BranchHash(branch string) string {
	log.Debugf("Getting hash for " + branch)
	output, err := execute("rev-parse", branch)
	if err != nil {
		log.Errorf("Error getting branch hash %s", err)
		os.Exit(1)
	}
	return output[0]
}
