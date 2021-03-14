package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	git "github.com/libgit2/git2go/v30"
	"github.com/stretchr/testify/require"
)

const (
	GIT_USERNAME = "Alice"
	GIT_EMAIL    = "alice@ecorp.com"
)

// TestRepo creates a new repository in a temporary directory
func TestRepo(t *testing.T) (repo *git.Repository) {
	path, err := ioutil.TempDir("", "turbogit")
	require.NoError(t, err)
	r, err := git.InitRepository(path, false)
	require.NoError(t, err)
	return r
}

// CleanupRepo removes the repo directory recursively
func CleanupRepo(t *testing.T, r *git.Repository) {
	p := r.Workdir()
	if r.IsBare() {
		p = r.Path()
	}
	require.NoError(t, os.RemoveAll(p))
}

// NewFile creates a new temp file in the repo working directory
func NewFile(t *testing.T, r *git.Repository) *os.File {
	f, err := ioutil.TempFile(r.Workdir(), "")
	require.NoError(t, err)
	return f
}

// StageFile adds a file to the repository index
func StageFile(t *testing.T, f *os.File, r *git.Repository) {
	frel, err := filepath.Rel(r.Workdir(), f.Name())
	require.NoError(t, err)
	idx, err := r.Index()
	require.NoError(t, err)
	require.NoError(t, idx.AddByPath(frel))
}

// StageNewFile creates a new file and adds it to the index
func StageNewFile(t *testing.T, r *git.Repository) {
	StageFile(t, NewFile(t, r), r)
}

// InitRepoConf set the repository initial configuration
func InitRepoConf(t *testing.T, r *git.Repository) {
	c, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, c.SetString("user.name", GIT_USERNAME))
	require.NoError(t, c.SetString("user.email", GIT_EMAIL))
}
