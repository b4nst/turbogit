package test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func CaptureStd(t *testing.T, std *os.File) (f *os.File, reset func()) {
	f, err := ioutil.TempFile("", path.Base(std.Name()))
	require.NoError(t, err)

	backup := *std
	reset = func() {
		*std = backup
	}
	*std = *f
	return
}

func WriteGitHook(t *testing.T, hook string, content string) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	hooks := path.Join(wd, ".git", "hooks")
	err = os.MkdirAll(hooks, 0700)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(hooks, hook), []byte(content), 0777)
	require.NoError(t, err)
}
