package git

import (
	"os"
	"strings"

	git2go "github.com/libgit2/git2go/v33"
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

func CurrentPatch(r *git2go.Repository) (string, error) {
	ref, err := r.RevparseSingle("HEAD")
	if err != nil {
		return "", err
	}
	old, err := ref.AsCommit()
	if err != nil {
		return "", err
	}

	tree, err := old.Tree()
	if err != nil {
		return "", err
	}
	diff, err := r.DiffTreeToIndex(tree, nil, &git2go.DiffOptions{
		Flags:            git2go.DiffIgnoreWhitespace,
		IgnoreSubmodules: git2go.SubmoduleIgnoreAll,
	})
	if err != nil {
		return "", err
	}
	numDeltas, err := diff.NumDeltas()
	if err != nil {
		return "", err
	}

	patches := make([]string, numDeltas)
	for i := 0; i < numDeltas; i++ {
		p, err := diff.Patch(i)
		if err != nil {
			return "", err
		}
		ps, err := p.String()
		if err != nil {
			return "", err
		}
		patches = append(patches, ps)
	}

	return strings.Join(patches, "\n"), nil
}
