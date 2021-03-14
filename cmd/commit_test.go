package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/b4nst/turbogit/internal/format"
	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommitCmd(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-commit")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)

	cmd := &cobra.Command{}

	cmd.Flags().StringP("type", "t", "fix", "")
	cmd.Flags().BoolP("breaking-changes", "c", true, "")
	cmd.Flags().BoolP("edit", "e", true, "")
	cmd.Flags().StringP("scope", "s", "scope", "")
	cmd.Flags().BoolP("amend", "a", true, "")

	cco, err := parseCommitCmd(cmd, []string{"hello", "world!"})
	require.NoError(t, err)
	expected := CommitCmdOption{
		CType:           format.FixCommit,
		Message:         "hello world!",
		Scope:           "scope",
		BreakingChanges: true,
		PromptEditor:    true,
		Amend:           true,
		Repo:            r,
	}
	assert.Equal(t, expected, *cco)
}
