package cmd

import (
	"testing"

	"github.com/b4nst/turbogit/internal/test"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchCreate(t *testing.T) {
	r, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}

	assertBranchIs := func(expected string) {
		ref, err := r.Head()
		require.NoError(t, err)
		assert.True(t, ref.Name().IsBranch(), "Should be a branch")
		assert.Equal(t, expected, ref.Name().Short())
	}

	// feature branch
	err = branch(cmd, []string{"feat"})
	assert.Error(t, err)
	err = branch(cmd, []string{"feat", "my", "feature"})
	assert.NoError(t, err)
	assertBranchIs("feat/my-feature")

	// fix branc
	err = branch(cmd, []string{"fix"})
	assert.Error(t, err)
	err = branch(cmd, []string{"fix", "my", "fix"})
	assert.NoError(t, err)
	assertBranchIs("fix/my-fix")

	// user with no description
	cfg := config.NewConfig()
	cfg.User = struct {
		Name  string
		Email string
	}{"bob", "bob@company.com"}
	r.SetConfig(cfg)
	err = branch(cmd, []string{"user"})
	assert.NoError(t, err)
	assertBranchIs("user/bob")

	// user with description
	cfg = config.NewConfig()
	cfg.User = struct {
		Name  string
		Email string
	}{"alice", "alice@company.com"}
	r.SetConfig(cfg)
	err = branch(cmd, []string{"user", "my", "branch"})
	assert.NoError(t, err)
	assertBranchIs("user/alice/my-branch")

}
