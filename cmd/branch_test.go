package cmd

import (
	"testing"

	"github.com/b4nst/turbogit/internal/format"
	tugit "github.com/b4nst/turbogit/internal/git"
	"github.com/b4nst/turbogit/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchCreate(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

	bco := &BranchCmdOption{
		format.TugBranch{Type: "feat", Description: "foo"},
		r,
	}
	assert.Error(t, runBranch(bco), "No commit to create branch from, please create the initial commit")
	tugit.Commit(r, "initial commit")

	assert.NoError(t, runBranch(bco))
	h, err := r.Head()
	require.NoError(t, err)
	assert.Equal(t, "refs/heads/feat/foo", h.Name())
}

func TestParseBranchCmd(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	test.InitRepoConf(t, r)

	cmd := &cobra.Command{}

	// User branch
	bco, err := parseBranchCmd(cmd, []string{"user", "my", "branch"})
	assert.NoError(t, err)
	expected := BranchCmdOption{
		NewBranch: format.TugBranch{Type: "user", Prefix: test.GIT_USERNAME, Description: "my branch"},
		Repo:      r,
	}
	assert.Equal(t, expected, *bco)
	// Classic branch
	bco, err = parseBranchCmd(cmd, []string{"feat", "foo", "bar"})
	assert.NoError(t, err)
	expected = BranchCmdOption{
		NewBranch: format.TugBranch{Type: "feat", Description: "foo bar"},
		Repo:      r,
	}
	assert.Equal(t, expected, *bco)
}
