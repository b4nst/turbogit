package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/b4nst/turbogit/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckCommit(t *testing.T) {
	r, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	c := test.AddCommit(t, r, "bad format")
	stderr, reset := test.CaptureStd(t, os.Stderr)
	ccr := checkCommit(c)
	reset()
	assert.False(t, ccr)
	stde, err := ioutil.ReadFile(stderr.Name())
	assert.Equal(t, fmt.Sprintf("%s %s\n", c.Hash, "bad format"), string(stde))

	c = test.AddCommit(t, r, format.CommitMessage(&format.CommitMessageOption{Ctype: format.BuildCommit, Description: "correct format"}))
	ccr = checkCommit(c)
	assert.True(t, ccr)
}

func TestRunCheck(t *testing.T) {
	r, teardown, err := test.SetUpRepo()
	defer teardown()
	require.NoError(t, err)

	cmd := &cobra.Command{}

	fAll := cmd.Flags().BoolP("all", "a", false, "Check all the refs in refs/, along with HEAD")
	fFrom := cmd.Flags().StringP("from", "f", "HEAD", "Hash of the commit to start from")

	// Add 3 commit (the first one is alread set by SetUpRepo)
	from := test.AddCommit(t, r, format.CommitMessage(&format.CommitMessageOption{Ctype: format.BuildCommit, Description: "hello world"}))
	test.AddCommit(t, r, format.CommitMessage(&format.CommitMessageOption{Ctype: format.FeatureCommit, Description: "right format"}))
	test.AddCommit(t, r, format.CommitMessage(&format.CommitMessageOption{Ctype: format.CiCommit, Description: "correct format"}))
	// Test with default option
	*fAll = false
	*fFrom = "HEAD"
	stdout, reset := test.CaptureStd(t, os.Stdout)
	err = runCheck(cmd, []string{})
	reset()
	assert.NoError(t, err)
	stdo, err := ioutil.ReadFile(stdout.Name())
	assert.Equal(t, "4 commit(s) checked.", string(stdo))
	// Test from
	*fAll = false
	*fFrom = from.Hash.String()
	stdout, reset = test.CaptureStd(t, os.Stdout)
	err = runCheck(cmd, []string{})
	reset()
	assert.NoError(t, err)
	stdo, err = ioutil.ReadFile(stdout.Name())
	assert.Equal(t, "2 commit(s) checked.", string(stdo))
	// Test from relative
	*fAll = false
	*fFrom = "HEAD~2"
	stdout, reset = test.CaptureStd(t, os.Stdout)
	err = runCheck(cmd, []string{})
	reset()
	assert.NoError(t, err)
	stdo, err = ioutil.ReadFile(stdout.Name())
	assert.Equal(t, "2 commit(s) checked.", string(stdo))
	// Test from bad
	*fAll = false
	*fFrom = "nope"
	err = runCheck(cmd, []string{})
	assert.EqualError(t, err, "Error looking for commit nope: reference not found")
	// Test with bad commit
	bad := test.AddCommit(t, r, "badbeef")
	*fAll = false
	*fFrom = "HEAD"
	stderr, reset := test.CaptureStd(t, os.Stderr)
	err = runCheck(cmd, []string{})
	reset()
	assert.EqualError(t, err, "The previous commit(s) do(es) not respect Conventional Commit.")
	stde, err := ioutil.ReadFile(stderr.Name())
	assert.Equal(t, fmt.Sprintf("%s %s\n", bad.Hash, "badbeef"), string(stde))
}
