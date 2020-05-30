package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestNeedCommit(t *testing.T) {
	r, teardown, err := setUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	nc, err := needCommit(ctx)
	assert.NoError(t, err)
	assert.False(t, nc)

	wd, err := os.Getwd()
	filename := filepath.Join(wd, "TestIsWorkingTreeClean")
	require.NoError(t, ioutil.WriteFile(filename, []byte("hello world!"), 0644))

	nc, err = needCommit(ctx)
	assert.EqualError(t, err, "no changes added to commit")

	wt, err := r.Worktree()
	require.NoError(t, err)
	_, err = wt.Add("TestIsWorkingTreeClean")
	assert.NoError(t, err)
	nc, err = needCommit(ctx)
	assert.NoError(t, err)
	assert.True(t, nc)
}

func stageNewFile(r *git.Repository) error {
	// Create and stage file
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile(wd, "*")
	if err != nil {
		return err
	}
	wt, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = wt.Add(filepath.Base(f.Name()))
	if err != nil {
		return err
	}
	return nil
}

func TestCommit(t *testing.T) {
	r, teardown, err := setUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}

	stdout := os.Stdout
	restore := func() {
		os.Stdout = stdout
	}
	defer restore()

	devnull, err := ioutil.TempFile("", "dev-null")
	require.NoError(t, err)
	defer os.RemoveAll(devnull.Name())
	os.Stdout = devnull

	fType := cmd.Flags().StringP("type", "t", "", "")
	fBreak := cmd.Flags().BoolP("breaking-changes", "c", false, "")
	cmd.Flags().BoolP("edit", "e", false, "")
	fScope := cmd.Flags().StringP("scope", "s", "", "")

	assertLastCommit := func(msg string) {
		citer, err := r.Log(&git.LogOptions{})
		require.NoError(t, err)
		c, err := citer.Next()
		require.NoError(t, err)
		assert.Equal(t, msg, c.Message)
	}

	// Bad commit type
	*fType = ""
	*fBreak = false
	*fScope = ""
	require.NoError(t, stageNewFile(r))
	assert.Error(t, commit(cmd, []string{"not-type"}))
	// Feat
	*fType = "feat"
	*fBreak = false
	*fScope = ""
	require.NoError(t, stageNewFile(r))
	assert.NoError(t, commit(cmd, []string{"my", "message"}))
	assertLastCommit("feat: my message")
	// Breaking change
	*fType = ""
	*fBreak = true
	*fScope = ""
	require.NoError(t, stageNewFile(r))
	assert.NoError(t, commit(cmd, []string{"fix", "my", "message"}))
	assertLastCommit("fix!: my message")
	// Scope
	*fType = ""
	*fBreak = false
	*fScope = "scope"
	require.NoError(t, stageNewFile(r))
	assert.NoError(t, commit(cmd, []string{"test", "my", "message"}))
	assertLastCommit("test(scope): my message")
	// Workdir clean
	*fType = ""
	*fBreak = false
	*fScope = ""
	assert.NoError(t, commit(cmd, []string{"fix", "not", "committed"}))
	assertLastCommit("test(scope): my message")
	// Unstaged files
	wd, err := os.Getwd()
	require.NoError(t, err)
	_, err = ioutil.TempFile(wd, "*")
	require.NoError(t, err)
	*fType = ""
	*fBreak = false
	*fScope = ""
	assert.EqualError(t, commit(cmd, []string{"fix", "not", "committed"}), "no changes added to commit")
}