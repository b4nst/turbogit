package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompletion(t *testing.T) {
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

	err = completion(cmd, []string{"bash"})
	assert.NoError(t, err)
	err = completion(cmd, []string{"zsh"})
	assert.NoError(t, err)
	err = completion(cmd, []string{"fish"})
	assert.NoError(t, err)
	err = completion(cmd, []string{"powershell"})
	assert.NoError(t, err)
	err = completion(cmd, []string{"other"})
	assert.Error(t, err)
}
