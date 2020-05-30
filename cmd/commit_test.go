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
