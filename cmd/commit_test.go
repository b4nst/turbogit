package cmd

import (
	"os"
	"testing"

	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	"github.com/b4nst/turbogit/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommitCmd(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	require.NoError(t, os.Chdir(r.Workdir()))

	cmd := &cobra.Command{}
	cmd.Flags().StringP("type", "t", "fix", "")
	cmd.Flags().BoolP("breaking-changes", "c", true, "")
	cmd.Flags().BoolP("edit", "e", true, "")
	cmd.Flags().StringP("scope", "s", "scope", "")
	cmd.Flags().BoolP("amend", "a", true, "")

	cmdbuilder.MockRepoAware(cmd, r)

	cco, err := parseCommitCmd(cmd, []string{"hello", "world!"})
	require.NoError(t, err)
	expect := commitOpt{
		CType:           format.FixCommit,
		Message:         "hello world!",
		Scope:           "scope",
		BreakingChanges: true,
		PromptEditor:    true,
		Amend:           true,
		Repo:            r,
	}
	assert.Equal(t, expect, *cco)
}
