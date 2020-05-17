package context

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func currRepo() (*git.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("not a git repository (or any of the parent directories): .git")
	}

	return r, nil
}
