package git

import (
	"errors"

	git "github.com/libgit2/git2go/v31"
)

// StageReady returns true if the stage is ready to be committed. Otherwise it returns false if there is nothing to commit or an error.
func StageReady(r *git.Repository) (bool, error) {
	s, err := r.StatusList(&git.StatusOptions{Show: git.StatusShowIndexAndWorkdir, Flags: git.StatusOptIncludeUntracked})
	if err != nil {
		return false, err
	}

	count, err := s.EntryCount()
	if err != nil {
		return false, err
	}
	if count <= 0 {
		return false, nil
	}
	for i := 0; i < count; i++ {
		se, err := s.ByIndex(i)
		if err != nil {
			return false, err
		}
		if se.Status <= git.StatusIndexTypeChange {
			return true, nil
		}
	}
	return false, errors.New("No changes added to commit")
}
