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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().StringP("type", "t", "f", "Branch type ([f]eature,[p]atch or [u]ser).")
}

var branchCmd = &cobra.Command{
	Use:           "branch [description]",
	Short:         "Create a new branch.",
	SilenceUsage:  false,
	SilenceErrors: true,
	Args:          validateArgs,
	RunE:          branch,
}

func branch(cmd *cobra.Command, args []string) error {
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}

	btype, err := validateBranchType(cmd)
	if err != nil {
		return err
	}

	headRef, err := ctx.Repo.Head()
	if err != nil {
		return err
	}

	// Format branch name
	if btype == format.UserBranch && ctx.Username == "" {
		return errors.New("You need to configure your username before creating a user branch.")
	}
	d := strings.Join(args, "-")
	branch_name := plumbing.NewBranchReferenceName(format.BranchName(btype, d, ctx.Username))

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

func validateArgs(cmd *cobra.Command, args []string) error {
	btype, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	if btype != "u" && len(args) < 1 {
		return fmt.Errorf("Type %s requires a description.", btype)
	}
	return nil
}

func validateBranchType(cmd *cobra.Command) (format.BranchType, error) {
	btype, err := cmd.Flags().GetString("type")
	if err != nil {
		return -1, err
	}

	switch btype {
	case "f":
		return format.FeatureBranch, nil
	case "p":
		return format.FixBranch, nil
	case "u":
		return format.UserBranch, nil
	default:
		return -1, fmt.Errorf("%s is not a branch type, allowed values are f,p and u", btype)
	}
}
