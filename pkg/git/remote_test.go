package git

import (
	"testing"

	"github.com/b4nst/turbogit/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRemote(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	_, err := r.Remotes.Create("origin", "git@alice.com:namespace/project.git")
	require.NoError(t, err)

	// Direct
	u, err := ParseRemote(r, "origin", false)
	assert.NoError(t, err)
	assert.Equal(t, u.String(), "ssh://git@alice.com/namespace/project.git")

	// No fallback
	u, err = ParseRemote(r, "rename", false)
	assert.EqualError(t, err, "remote 'rename' does not exist")

	// Fallback
	u, err = ParseRemote(r, "rename", true)
	assert.NoError(t, err)
	assert.Equal(t, u.String(), "file://origin")

}
