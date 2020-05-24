package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastTag(t *testing.T) {
	r, teardown, err := setUpRepo()
	defer teardown()
	require.NoError(t, err)

	// No tag
	last, err := getLastTag(r)
	require.NoError(t, err)
	assert.Nil(t, last)

	testNewVersion := func(v string) {
		head, err := r.Head()
		require.NoError(t, err)
		ref, err := r.CreateTag(v, head.Hash(), nil)
		require.NoError(t, err)
		err = r.Storer.SetReference(ref)
		require.NoError(t, err)
		last, err = getLastTag(r)
		require.NoError(t, err)
		require.NotNil(t, last)
		assert.Equal(t, ref, last.ref)
		assert.Equal(t, v, last.version.String())
	}

	testNewVersion("1.0.0-alpha")
	testNewVersion("1.0.0-alpha.1")
	testNewVersion("1.0.0-beta")
	testNewVersion("1.0.0-beta.1")
	testNewVersion("1.0.1")
	testNewVersion("1.1.0")
	testNewVersion("2.0.0")
}
