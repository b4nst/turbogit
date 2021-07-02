package git

import (
	"os"

	git2go "github.com/libgit2/git2go/v31"
)

// Getrepo returns the repository in the current directory or an error.
func Getrepo() (*git2go.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	rpath, err := git2go.Discover(wd, false, nil)
	if err != nil {
		return nil, err
	}
	repo, err := git2go.OpenRepository(rpath)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
