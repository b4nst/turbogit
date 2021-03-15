package cmd

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/b4nst/turbogit/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunVersion(t *testing.T) {
	f, restore := test.CaptureStd(t, os.Stdout)
	defer restore()
	runVersion(&cobra.Command{}, []string{})

	content, err := ioutil.ReadFile(f.Name())
	require.NoError(t, err)
	assert.Contains(t, string(content), Version, "Version should contains turbogit version")
	assert.Contains(t, string(content), runtime.Version(), "Version should contains runtime version")
	assert.Contains(t, string(content), Commit, "Version should contains build commit")
	assert.Contains(t, string(content), BuildDate, "Version should contains build date")
	assert.Contains(t, string(content), runtime.GOARCH, "Version should contains go arch")
	assert.Contains(t, string(content), runtime.GOOS, "Version should contains go OS")
}
