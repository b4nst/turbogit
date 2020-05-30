package cmd

import (
	"testing"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteCommit(t *testing.T) {
	r, teardown, err := setUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)
	assert.NoError(t, writeCommit(ctx, "commit message"))

	citer, err := r.Log(&git.LogOptions{})
	require.NoError(t, err)

	c, err := citer.Next()
	require.NoError(t, err)
	assert.Equal(t, "commit message", c.Message)
}

func TestCommit(t *testing.T) {
	r, teardown, err := setUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}

	fType := cmd.Flags().StringP("type", "t", "", "")
	fBreak := cmd.Flags().BoolP("breaking-changes", "c", false, "")
	cmd.Flags().BoolP("edit", "e", false, "")
	fScope := cmd.Flags().StringP("scope", "s", "", "")

	assert.Error(t, commit(cmd, []string{"not-type"}))

	assertLastCommit := func(msg string) {
		citer, err := r.Log(&git.LogOptions{})
		require.NoError(t, err)
		c, err := citer.Next()
		require.NoError(t, err)
		assert.Equal(t, msg, c.Message)
	}

	*fType = "feat"
	*fBreak = false
	*fScope = ""
	assert.NoError(t, commit(cmd, []string{"my", "message"}))
	assertLastCommit("feat: my message")

	*fType = ""
	*fBreak = true
	*fScope = ""
	assert.NoError(t, commit(cmd, []string{"fix", "my", "message"}))
	assertLastCommit("fix!: my message")

	*fType = ""
	*fBreak = false
	*fScope = "scope"
	assert.NoError(t, commit(cmd, []string{"test", "my", "message"}))
	assertLastCommit("test(scope): my message")
}
