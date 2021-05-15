package git

import (
	"testing"

	"github.com/b4nst/turbogit/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestStageReady(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)

	nc, err := StageReady(r)
	assert.NoError(t, err)
	assert.False(t, nc)

	f := test.NewFile(t, r)
	nc, err = StageReady(r)
	assert.EqualError(t, err, "No changes added to commit")

	test.StageFile(t, f, r)
	nc, err = StageReady(r)
	assert.NoError(t, err)
	assert.True(t, nc)
}
