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
	"fmt"

	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	"github.com/hashicorp/go-multierror"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.Flags().BoolP("all", "a", false, "Check all the commits in refs/*, along with HEAD")
	RootCmd.Flags().StringP("from", "f", "HEAD", "Commit to start from. Can be a hash or any revision as accepted by rev parse.")

	cmdbuilder.RepoAware(RootCmd)
}

type option struct {
	All  bool
	From string
	Repo *git.Repository
}

var RootCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the history to follow conventional commit",
	Example: `
# Check if all is ok
$ git check
`,
	Args: cobra.NoArgs,
	Run:  run,
}

func run(cmd *cobra.Command, args []string) {
	opt := &option{}
	var err error

	opt.All, err = cmd.Flags().GetBool("all")
	cobra.CheckErr(err)

	opt.From, err = cmd.Flags().GetString("from")
	cobra.CheckErr(err)

	opt.Repo = cmdbuilder.GetRepo(cmd)

	cobra.CheckErr(check(opt))

	cmd.Println("repository compliant.")
}

func check(opt *option) error {
	walk, err := opt.Repo.Walk()
	if err != nil {
		return err
	}
	if opt.All {
		if err := walk.PushGlob("refs/*"); err != nil {
			return err
		}
	} else {
		from, err := opt.Repo.RevparseSingle(opt.From)
		if err != nil {
			return err
		}
		if err := walk.Push(from.Id()); err != nil {
			return err
		}
	}

	merr := &multierror.Error{}
	if err := walk.Iterate(walker(merr)); err != nil {
		return err
	}
	return merr.ErrorOrNil()
}

func walker(merr *multierror.Error) git.RevWalkIterator {
	return func(c *git.Commit) bool {
		sid, err := c.ShortId()
		if err != nil {
			multierror.Append(merr, err)
			return true
		}
		co := format.ParseCommitMsg(c.Message())
		if co == nil {
			multierror.Append(merr, fmt.Errorf("%s ('%s') is not compliant", sid, c.Summary()))
		}
		return true
	}
}
