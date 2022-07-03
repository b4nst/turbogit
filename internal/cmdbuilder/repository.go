package cmdbuilder

import (
	"context"

	tugit "github.com/b4nst/turbogit/pkg/git"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

type repoKey struct{}

// GetRepo returns the current git repository.
func GetRepo(cmd *cobra.Command) *git.Repository {
	v := cmd.Context().Value(repoKey{})
	return v.(*git.Repository)
}

func RepoAware(cmd *cobra.Command) {
	AppendPreRun(cmd, repoPreRun)
}

func MockRepoAware(cmd *cobra.Command, repo *git.Repository) {
	parent := cmd.Context()
	if parent == nil {
		parent = context.TODO()
	}
	cmd.SetContext(context.WithValue(parent, repoKey{}, repo))
}

func repoPreRun(cmd *cobra.Command, args []string) {
	repo, err := tugit.Getrepo()
	cobra.CheckErr(err)

	cmd.SetContext(context.WithValue(cmd.Context(), repoKey{}, repo))
}
