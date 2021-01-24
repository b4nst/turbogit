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
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/b4nst/turbogit/internal/format"

	"github.com/libgit2/git2go/v30"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(branchCmd)
}

type BranchCmdOption struct {
	BType format.BranchType
	Name  string
	Repo  *git.Repository
}

var branchCmd = &cobra.Command{
	Use:                   fmt.Sprintf("branch %s [description]", format.AllBranchType()),
	Short:                 "Create a new branch",
	DisableFlagsInUseLine: true,
	Example: `
# Create branch feat/my-feature from current branch
$ tug branch feat my feature

# Create branch user/alice/my-branch, given that alice is the current tug/git user
$ tug branch user my branch
	`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	ValidArgs:    format.AllBranchType(),
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
	btype, err := format.BranchTypeFrom(args[0])
	if err != nil {
		return nil, err
	}
	if btype != format.UserBranch && len(args) < 2 {
		return nil, fmt.Errorf("%s branches need a description.", btype)
	}

	// Find repo
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	rpath, err := git.Discover(wd, false, nil)
	if err != nil {
		return nil, err
	}
	repo, err := git.OpenRepository(rpath)
	if err != nil {
		return nil, err
	}

	// Get user name from config
	config, err := repo.Config()
	if err != nil {
		return nil, err
	}
	username, _ := config.LookupString("user.name")
	// Format branch name
	if btype == format.UserBranch && username == "" {
		return nil, errors.New("You need to configure your username before creating a user branch.")
	}

	return &BranchCmdOption{
		BType: btype,
		Name:  format.BranchName(btype, strings.Join(args[1:], "-"), username),
		Repo:  repo,
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
	b, err := r.CreateBranch(bco.Name, t, false)
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
