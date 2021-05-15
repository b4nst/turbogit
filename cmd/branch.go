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
	"errors"
	"log"
	"strings"

	"github.com/b4nst/turbogit/pkg/format"
	tugit "github.com/b4nst/turbogit/pkg/git"
	"github.com/b4nst/turbogit/pkg/integrations"

	git "github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(branchCmd)
}

type BranchCmdOption struct {
	NewBranch format.TugBranch
	Repo      *git.Repository
}

var branchCmd = &cobra.Command{
	Use:                   "branch [type] [description]",
	Aliases:               []string{"b"},
	Short:                 "Create a new branch",
	DisableFlagsInUseLine: true,
	Long: `
If you don't give any argument, the command will look for issue in pre-configured issues provider.
The issue ID will be used as a prefix.
If type=user, a prefix with your git username will be added to the branch name.
	`,
	Example: `
# Create branch feat/my-feature from current branch
$ tug branch feat my feature

# Create branch user/alice/my-branch, given that alice is your git username
$ tug branch user my branch

# Look for an issue provider, prompt for issue and create a branch from there
$ tug branch
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return errors.New("accepts 0 or at least 2 args, received 1")
		}
		return nil
	},
	SilenceUsage: true,
	Run:          runBranchCmd,
}

func runBranchCmd(cmd *cobra.Command, args []string) {
	bco, err := parseBranchCmd(cmd, args)
	if err != nil {
		log.Fatal(err)
	}
	err = runBranch(bco)
	if err != nil {
		log.Fatal(err)
	}
}

func parseBranchCmd(cmd *cobra.Command, args []string) (*BranchCmdOption, error) {
	// Find repo
	repo, err := tugit.Getrepo()
	if err != nil {
		return nil, err
	}

	var nb format.TugBranch

	if len(args) == 0 {
		providers, err := integrations.ProvidersFrom(repo)
		if err != nil {
			return nil, err
		}
		issues := []integrations.IssueDescription{}
		for _, p := range providers {
			// TODO concurrent search
			pIssues, err := p.Search()
			if err != nil {
				return nil, err
			}
			issues = append(issues, pIssues...)
		}
		issue, err := integrations.SelectIssue(issues, false)
		if err != nil {
			return nil, err
		}
		nb = issue.ToBranch(format.DefaultTypeRewrite)
	} else {
		nb = format.TugBranch{
			Description: strings.Join(args[1:], " "),
		}.WithType(args[0], format.DefaultTypeRewrite)
		if nb.Type == "user" {
			// Get user name from config
			config, err := repo.Config()
			if err != nil {
				return nil, err
			}
			username, _ := config.LookupString("user.name")
			if username == "" {
				return nil, errors.New("You need to configure your username before creating a user branch.")
			}
			nb.Prefix = username
		}
	}

	return &BranchCmdOption{
		NewBranch: nb,
		Repo:      repo,
	}, nil
}

func runBranch(bco *BranchCmdOption) error {
	r := bco.Repo

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
	b, err := r.CreateBranch(bco.NewBranch.String(), t, false)
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
