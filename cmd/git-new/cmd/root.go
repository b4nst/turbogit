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
	"strings"

	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	"github.com/b4nst/turbogit/pkg/integrations"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	cmdbuilder.RepoAware(RootCmd)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "new [type] [description]",
	Short: "Start a new branch.",
	Long: `
If you don't give any argument, the command will look for issue in pre-configured issues provider.
The issue ID will be used as a prefix.
If type=user(s), a prefix with your git username will be added to the branch name.
	`,
	Example: `
# Start new feature feat/my-feature from current branch
$ git new feat my feature

# Start working on a user branch (my-branch). This will create user/alice/my-branch, given that alice is your git username
$ tug branch user my branch

# Start working on a new issue (an issue provider must be configured on the repositoty)
$ tug branch
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return errors.New("accepts 0 or at least 2 args, received 1")
		}
		return nil
	},
	Run: run,
}

type option struct {
	NewBranch format.TugBranch
	Repo      *git.Repository
}

func run(cmd *cobra.Command, args []string) {
	opt := &option{}
	var err error

	opt.Repo = cmdbuilder.GetRepo(cmd)

	if len(args) <= 0 {
		opt.NewBranch, err = promptProviderBranch(opt.Repo)
		cobra.CheckErr(err)
	} else {
		opt.NewBranch = format.TugBranch{Description: strings.Join(args[1:], " ")}.
			WithType(args[0], format.DefaultTypeRewrite)
	}

	// User(s) branch
	if opt.NewBranch.Type == "user" || opt.NewBranch.Type == "users" {
		// Get user name from config
		config, err := opt.Repo.Config()
		cobra.CheckErr(err)

		username, _ := config.LookupString("user.name")
		if username == "" {
			cobra.CheckErr("You need to configure your username before creating a user branch.")
		}
		opt.NewBranch.Prefix = username
	}

	cobra.CheckErr(gnew(opt))
}

func gnew(opt *option) error {
	r := opt.Repo

	var t *git.Commit = nil
	head, err := r.Head()
	if err == nil {
		t, err = r.LookupCommit(head.Target())
		if err != nil {
			return err
		}
	}
	if t == nil {
		return errors.New("No commit to create branch from, please create the initial commit")
	}

	// Create new branch
	b, err := r.CreateBranch(opt.NewBranch.String(), t, false)
	if err != nil {
		return err
	}
	bc, err := r.LookupCommit(b.Target())
	if err != nil {
		return err
	}
	tree, err := r.LookupTree(bc.TreeId())
	if err != nil {
		return err
	}

	// Checkout the branch
	err = r.CheckoutTree(tree, &git.CheckoutOpts{Strategy: git.CheckoutSafe})
	if err != nil {
		return err
	}
	err = r.SetHead(b.Reference.Name())
	return err
}

func promptProviderBranch(repo *git.Repository) (nb format.TugBranch, err error) {
	providers, err := integrations.ProvidersFrom(repo)
	if err != nil {
		return nb, err
	}
	issues := []integrations.IssueDescription{}
	for _, p := range providers {
		// TODO concurrent search
		pIssues, err := p.Search()
		if err != nil {
			return nb, err
		}
		issues = append(issues, pIssues...)
	}
	issue, err := integrations.SelectIssue(issues, false)
	if err != nil {
		return nb, err
	}

	return issue.ToBranch(format.DefaultTypeRewrite), nil
}
