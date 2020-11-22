package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/b4nst/turbogit/internal/format"
	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteCommit(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-commit")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	config, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, config.SetString("user.name", "alice"))
	require.NoError(t, config.SetString("user.email", "alice@ecorp.com"))
	f, err := ioutil.TempFile(dir, "test-commit")
	require.NoError(t, err)
	frel, err := filepath.Rel(dir, f.Name())
	require.NoError(t, err)
	idx, err := r.Index()
	require.NoError(t, err)
	require.NoError(t, idx.AddByPath(frel))

	commit, err := writeCommit(r, "commit message")
	assert.NoError(t, err)
	assert.Equal(t, "commit message", commit.Message())
	assert.Equal(t, "alice", commit.Author().Name)
	assert.Equal(t, "alice@ecorp.com", commit.Author().Email)
	assert.WithinDuration(t, time.Now(), commit.Author().When, time.Second)
	head, err := r.Head()
	require.NoError(t, err)
	headCommit, err := r.LookupCommit(head.Target())
	require.NoError(t, err)
	assert.Equal(t, headCommit.Id(), commit.Id())
}

func TestNeedCommit(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-commit")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)

	nc, err := needCommit(r)
	assert.NoError(t, err)
	assert.False(t, nc)

	fmt.Println("Writing file")
	filename := filepath.Join(dir, "TestIsWorkingTreeClean")
	require.NoError(t, ioutil.WriteFile(filename, []byte("hello world!"), 0644))
	fmt.Println("Writing file done")

	nc, err = needCommit(r)
	assert.EqualError(t, err, "No changes added to commit")

	idx, err := r.Index()
	require.NoError(t, err)
	require.NoError(t, idx.AddByPath("TestIsWorkingTreeClean"))

	nc, err = needCommit(r)
	assert.NoError(t, err)
	assert.True(t, nc)
}

func TestSignature(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-commit")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	config, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, config.SetString("user.name", "alice"))
	require.NoError(t, config.SetString("user.email", "alice@ecorp.com"))

	sig, err := signature(r)
	require.NoError(t, err)
	assert.Equal(t, "alice", sig.Name)
	assert.Equal(t, "alice@ecorp.com", sig.Email)
	assert.WithinDuration(t, time.Now(), sig.When, time.Second)
}

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

	cco, err := parseCommitCmd(cmd, []string{"hello", "world!"})
	require.NoError(t, err)
	expected := CommitCmdOption{
		CType:           format.FixCommit,
		Message:         "hello world!",
		Scope:           "scope",
		BreakingChanges: true,
		PromptEditor:    true,
		Repo:            r,
	}
	assert.Equal(t, expected, *cco)
}
