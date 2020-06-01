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
	"strings"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/format"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
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
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: format.AllBranchType(),
	RunE:      branch,
}

func branch(cmd *cobra.Command, args []string) error {
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}

	btype, err := format.BranchTypeFrom(args[0])
	if err != nil {
		return err
	}
	if btype != format.UserBranch && len(args) < 2 {
		return fmt.Errorf("%s branches need a description.", btype)
	}

	headRef, err := ctx.Repo.Head()
	if err != nil {
		return err
	}

	// Format branch name
	cfg, err := ctx.Repo.ConfigScoped(config.SystemScope)
	if btype == format.UserBranch && cfg.User.Name == "" {
		return errors.New("You need to configure your username before creating a user branch.")
	}
	d := strings.Join(args[1:], "-")
	branch_name := plumbing.NewBranchReferenceName(format.BranchName(btype, d, cfg.User.Name))

	// Create new branch
	ref := plumbing.NewHashReference(branch_name, headRef.Hash())
	if err := ctx.Repo.Storer.SetReference(ref); err != nil {
		return err
	}

	// Checkout former branch
	w, err := ctx.Repo.Worktree()
	if err != nil {
		return err
	}
	if err := w.Checkout(&git.CheckoutOptions{Branch: ref.Name()}); err != nil {
		return err
	}

	return nil
}
