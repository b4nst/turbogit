/*
Copyright © 2020 banst

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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/format"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:                   "check",
	Short:                 "Check the history to follow conventional commit",
	DisableFlagsInUseLine: true,
	Example: `
# Check if all is ok
$ tug check
`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE:         runCheck,
}

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("all", "a", false, "Check all the refs in refs/, along with HEAD")
	checkCmd.Flags().StringP("from", "f", "HEAD", "Hash of the commit to start from")
}

func runCheck(cmd *cobra.Command, args []string) error {
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}
	// Flags
	// --all
	fAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}
	// --from
	fFrom, err := cmd.Flags().GetString("from")
	if err != nil {
		return err
	}
	from, err := ctx.Repo.ResolveRevision(plumbing.Revision(fFrom))
	if err != nil {
		return fmt.Errorf("Error looking for commit %s: %s", fFrom, err)
	}

	cIter, err := ctx.Repo.Log(&git.LogOptions{Order: git.LogOrderCommitterTime, All: fAll, From: *from})
	clean, checked := true, 0
	for c, err := cIter.Next(); err == nil; c, err = cIter.Next() {
		clean, checked = checkCommit(c) && clean, checked+1
	}
	if err != nil && err != io.EOF {
		return err
	}
	if !clean {
		return fmt.Errorf("The previous commit(s) do(es) not respect Conventional Commit.")
	}
	fmt.Printf("%d commit(s) checked.", checked)
	return nil
}

func checkCommit(c *object.Commit) bool {
	co := format.ParseCommitMsg(c.Message)
	if co == nil {
		fmt.Fprintln(os.Stderr, c.Hash, strings.SplitN(c.Message, "\n", 2)[0])
		return false
	}
	return true
}
