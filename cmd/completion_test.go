package cmd

import (
	"os"
	"testing"

	"github.com/b4nst/turbogit/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	cmd := &cobra.Command{}

	f, restore := test.CaptureStd(t, os.Stdout)
	defer restore()
	defer os.RemoveAll(f.Name())

	assert.NoError(t, runCompletion(cmd, []string{"bash"}))
	assert.NoError(t, runCompletion(cmd, []string{"zsh"}))
	assert.NoError(t, runCompletion(cmd, []string{"fish"}))
	assert.NoError(t, runCompletion(cmd, []string{"powershell"}))
	assert.Error(t, runCompletion(cmd, []string{"other"}))
}
