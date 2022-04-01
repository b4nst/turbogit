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
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/b4nst/turbogit/pkg/format"
	tugit "github.com/b4nst/turbogit/pkg/git"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().BoolP("all", "a", false, "Check all the commits in refs/*, along with HEAD")
	rootCmd.Flags().StringP("from", "f", "HEAD", "Commit to start from. Can be a hash or any revision as accepted by rev parse.")
}

type option struct {
	All  bool
	From string
	Repo *git.Repository
}

var rootCmd = &cobra.Command{
	Use:   "git-check",
	Short: "Check the history to follow conventional commit",
	Example: `
# Check if all is ok
$ git check
`,
	Args: cobra.NoArgs,
	Run:  runCmd,
}

func runCmd(cmd *cobra.Command, args []string) {
	opt, err := parseCmd(cmd, args)
	if err != nil {
		log.Fatal(err)
	}
	err = run(opt)
	if err != nil {
		log.Fatal(err)
	}
}

func parseCmd(cmd *cobra.Command, args []string) (*option, error) {
	// --all
	fAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		return nil, err
	}
	// --from
	fFrom, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, err
	}

	// Find repo
	repo, err := tugit.Getrepo()
	if err != nil {
		return nil, err
	}

	return &option{
		All:  fAll,
		From: fFrom,
		Repo: repo,
	}, nil
}

func run(opt *option) error {
	r := opt.Repo

	walk, err := r.Walk()
	if err != nil {
		return err
	}
	if opt.All {
		if err := walk.PushGlob("refs/*"); err != nil {
			return err
		}
	} else {
		from, err := r.RevparseSingle(opt.From)
		if err != nil {
			return err
		}
		if err := walk.Push(from.Id()); err != nil {
			return err
		}
	}

	// Non format compliant commits
	var nfc []git.Commit

	walker := func(c *git.Commit) bool {
		co := format.ParseCommitMsg(c.Message())
		if co == nil {
			nfc = append(nfc, *c)
		}
		return true
	}
	if err := walk.Iterate(walker); err != nil {
		return err
	}
	if len(nfc) == 0 {
		fmt.Println("All commits are compliant")
		return nil
	} else {
		for _, c := range nfc {
			sid, err := c.ShortId()
			if err != nil {
				sid = c.Id().String()
			}
			fmt.Fprintln(os.Stderr, sid, c.Summary())
		}
		return errors.New("This commits are not compliant")
	}
}
