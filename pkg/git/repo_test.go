package git

import (
	"io/ioutil"
	"os"
	"testing"

	git "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetrepo(t *testing.T) {
	dir, err := ioutil.TempDir("", "turbogit-test-getrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, os.Chdir(dir))

	// Working dir is not a repo
	_, err = Getrepo()
	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "could not find repository from")
	}

	// Working dir is a repo
	r, err := git.InitRepository(dir, false)
	require.NoError(t, err)
	repo, err := Getrepo()
	assert.NoError(t, err)
	assert.Equal(t, r, repo)
}
