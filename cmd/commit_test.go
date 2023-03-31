/*
Copyright © 2022 banst

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"testing"

	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	"github.com/b4nst/turbogit/pkg/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommitCmd(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	require.NoError(t, os.Chdir(r.Workdir()))

	cmd := &cobra.Command{}
	cmd.Flags().StringP("type", "t", "fix", "")
	cmd.Flags().BoolP("breaking-changes", "c", true, "")
	cmd.Flags().BoolP("edit", "e", true, "")
	cmd.Flags().StringP("scope", "s", "scope", "")
	cmd.Flags().BoolP("amend", "a", true, "")
	cmd.Flags().BoolP("fill", "f", true, "")

	cmdbuilder.MockRepoAware(cmd, r)

	cco, err := parseCommitCmd(cmd, []string{"hello", "world!"})
	require.NoError(t, err)
	expect := commitOpt{
		CType:           format.FixCommit,
		Message:         "hello world!",
		Scope:           "scope",
		BreakingChanges: true,
		PromptEditor:    true,
		Amend:           true,
		Repo:            r,
		Fill:            true,
	}
	assert.Equal(t, expect, *cco)
}
