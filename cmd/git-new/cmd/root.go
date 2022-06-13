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
	"log"
	"strings"

	"github.com/b4nst/turbogit/pkg/format"
	"github.com/b4nst/turbogit/pkg/integrations"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "git-new [type] [description]",
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
	Run: runCmd,
}

type option struct {
	NewBranch format.TugBranch
	Repo      *git.Repository
}

func runCmd(cmd *cobra.Command, args []string) {
	bco, err := parseCmd(cmd, args)
	if err != nil {
		log.Fatal(err)
	}
	err = run(bco)
	if err != nil {
		log.Fatal(err)
	}
}

func parseCmd(cmd *cobra.Command, args []string) (*option, error) {
	// Find repo
	repo, err := git.PlainOpenWithOptions("", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, err
	}

	var nb format.TugBranch

	if len(args) == 0 {
		nb, err = promptProviderBranch(repo)
		if err != nil {
			return nil, err
		}
	} else {
		nb = format.TugBranch{
			Description: strings.Join(args[1:], " "),
		}.WithType(args[0], format.DefaultTypeRewrite)
	}

	// User(s) branch
	if nb.Type == "user" || nb.Type == "users" {
		// Get user name from config
		config, err := repo.Config()
		if err != nil {
			return nil, err
		}
		username := config.User.Name
		if username == "" {
			return nil, errors.New("You need to configure your username before creating a user branch.")
		}
		nb.Prefix = username
	}

	return &option{
		NewBranch: nb,
		Repo:      repo,
	}, nil
}

func run(opt *option) error {
	r := opt.Repo
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(opt.NewBranch.String()),
		Create: true,
		Keep:   true,
	})
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
