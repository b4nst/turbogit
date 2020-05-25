package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromCommand(t *testing.T) {
	expect, teardown, err := setUp()
	defer teardown()
	require.NoError(t, err)

	ctx, err := FromCommand(nil)
	require.NoError(t, err)
	assert.Equal(t, expect, ctx.Repo)
}
