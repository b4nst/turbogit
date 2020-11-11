package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHookCmd(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	// Test when no hooks exists
	hook := "hook-script"
	hc, err := hookCmd(ctx, hook)
	assert.NoError(t, err)
	assert.Nil(t, hc)

	// Test error with directory script instead of file
	err = os.MkdirAll(path.Join(".git", "hooks", hook), 0700)
	require.NoError(t, err)
	hc, err = hookCmd(ctx, hook)
	assert.EqualError(t, err, fmt.Sprintf("Hook .git/hooks/%s is a directory, it should be an executable file.", hook))
	assert.Nil(t, hc)
	err = os.Remove(path.Join(".git", "hooks", hook))
	require.NoError(t, err)

	// Test command
	test.WriteGitHook(t, hook, "")
	hc, err = hookCmd(ctx, hook)
	assert.NoError(t, err)
	cwd, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, &exec.Cmd{
		Dir:    cwd,
		Path:   path.Join(cwd, ".git", "hooks", hook),
		Args:   []string{path.Join(cwd, ".git", "hooks", hook)},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, hc)
}

func TestNoArgHook(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	hook := "hook-script"

	// Test without script
	err = noArgHook(ctx, hook)
	assert.NoError(t, err)

	// Test error script
	script := `#!/bin/sh
>&2 echo standard error
exit 3
`
	test.WriteGitHook(t, hook, script)
	stderr, resetSterr := test.CaptureStd(t, os.Stderr)
	defer resetSterr()
	err = noArgHook(ctx, hook)
	assert.EqualError(t, err, "exit status 3")
	stde, err := ioutil.ReadFile(stderr.Name())
	require.NoError(t, err)
	assert.Equal(t, "standard error\n", string(stde))

	// Test successful script
	script = `#!/bin/sh
echo Hello world!
exit 0
`
	test.WriteGitHook(t, hook, script)
	stdout, resetStdout := test.CaptureStd(t, os.Stdout)
	defer resetStdout()
	err = noArgHook(ctx, hook)
	assert.NoError(t, err)
	stdo, err := ioutil.ReadFile(stdout.Name())
	require.NoError(t, err)
	assert.Equal(t, "Hello world!\n", string(stdo))
}

func TestFileHook(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	hook := "hook-script"

	// Test without script
	msg, err := fileHook(ctx, hook, "hello world!")
	assert.NoError(t, err)
	assert.Equal(t, "hello world!", msg)

	// Test error script
	script := `#!/bin/sh
>&2 echo standard error
exit 3
`
	test.WriteGitHook(t, hook, script)
	stderr, resetSterr := test.CaptureStd(t, os.Stderr)
	defer resetSterr()
	msg, err = fileHook(ctx, hook, "hello world!")
	assert.EqualError(t, err, "exit status 3")
	assert.Equal(t, "hello world!", msg)
	stde, err := ioutil.ReadFile(stderr.Name())
	require.NoError(t, err)
	assert.Equal(t, "standard error\n", string(stde))

	// Test successful script
	script = `#!/bin/sh
echo "Hello world!" > "$1"
exit 0
`
	test.WriteGitHook(t, hook, script)
	msg, err = fileHook(ctx, hook, "Hey you!")
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!\n", msg)
}

func TestPreCommitHook(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	script := `#!/bin/sh
echo Hello world!
exit 0
`
	test.WriteGitHook(t, "pre-commit", script)
	stdout, resetStdout := test.CaptureStd(t, os.Stdout)
	defer resetStdout()
	err = PreCommitHook(ctx)
	assert.NoError(t, err)
	stdo, err := ioutil.ReadFile(stdout.Name())
	require.NoError(t, err)
	assert.Equal(t, "Hello world!\n", string(stdo))
}

func TestPrepareCommitMsg(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	// Test successful script
	script := `#!/bin/sh
echo "Hello world!" > "$1"
exit 0
`
	test.WriteGitHook(t, "prepare-commit-msg", script)
	msg, err := PrepareCommitMsgHook(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!\n", msg)
}

func TestCommitMsg(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	// Test successful script
	script := `#!/bin/sh
echo world! >> "$1"
exit 0
`
	test.WriteGitHook(t, "commit-msg", script)
	msg, err := CommitMsgHook(ctx, "Hello ")
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!\n", msg)
}

func TestPostCommit(t *testing.T) {
	_, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}
	ctx, err := context.FromCommand(cmd)
	require.NoError(t, err)

	script := `#!/bin/sh
echo Hello world!
exit 0
`
	test.WriteGitHook(t, "post-commit", script)
	stdout, resetStdout := test.CaptureStd(t, os.Stdout)
	defer resetStdout()
	err = PostCommitHook(ctx)
	assert.NoError(t, err)
	stdo, err := ioutil.ReadFile(stdout.Name())
	require.NoError(t, err)
	assert.Equal(t, "Hello world!\n", string(stdo))
}
