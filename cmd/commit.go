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

	"github.com/AlecAivazis/survey/v2"
	"github.com/b4nst/turbogit/internal/context"
	"github.com/b4nst/turbogit/internal/format"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("type", "t", "", fmt.Sprintf("Commit types %s", format.AllCommitType()))
	commitCmd.RegisterFlagCompletionFunc("type", typeFlagCompletion)
	commitCmd.Flags().BoolP("breaking-changes", "c", false, "Commit contains breaking changes")
	commitCmd.Flags().BoolP("edit", "e", false, "Prompt editor to edit your message (add body or/and footer(s))")
	commitCmd.Flags().StringP("scope", "s", "", "Add a scope")
}

func typeFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return format.AllCommitType(), cobra.ShellCompDirectiveDefault
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:                   "commit [type] [subject]",
	Short:                 "Commit staging area",
	DisableFlagsInUseLine: true,
	Example: `
# Commit a new feature (feat: a new feature)
$ tug commit feat a new feature

# Commit a fix that brings breaking changes (fix!: API break)
$ tug commit fix -c API break

# Add a scope to the commit (refactor(scope): a scopped refactor)
$ tug commit refactor a scopped refactor -s scope

# Open your editor to edit the commit message
$ tug commit ci -e message
	`,
	Args:      cobra.MinimumNArgs(1),
	RunE:      commit,
	ValidArgs: format.AllCommitType(),
}

func commit(cmd *cobra.Command, args []string) error {
	// Get context
	ctx, err := context.FromCommand(cmd)
	if err != nil {
		return err
	}
	// Parse flags
	typeFlag, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	if typeFlag == "" {
		typeFlag = args[0]
		args = args[1:]
	}
	bc, err := cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return err
	}
	scope, err := cmd.Flags().GetString("scope")
	if err != nil {
		return err
	}
	edit, err := cmd.Flags().GetBool("edit")
	if err != nil {
		return err
	}

	// Check if working tree is clean
	nc, err := needCommit(ctx)
	if err != nil {
		return err
	}
	if !nc {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	// Parse commit type
	ctype := format.FindCommitType(typeFlag)
	if ctype == format.NilCommit {
		return fmt.Errorf("'%s' is not a valid commit type", typeFlag)
	}
	// Create message
	msg := strings.Join(args, " ")
	cmsg := format.CommitMessage(&format.CommitMessageOption{
		Ctype: ctype, BreakingChanges: bc, Description: msg, Scope: scope,
	})
	if edit {
		cmsg = promptEditor(cmsg)
	}

	// Write commit
	if err := writeCommit(ctx, cmsg); err != nil {
		return err
	}

	return nil
}

func needCommit(ctx *context.Context) (bool, error) {
	wt, err := ctx.Repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := wt.Status()

	if err != nil {
		return false, err
	}

	unstaged := false
	for _, s := range status {
		if s.Staging != git.Unmodified && s.Staging != git.Untracked {
			return true, nil
		}
		if s.Worktree != git.Unmodified {
			unstaged = true
		}
	}

	if unstaged {
		return false, errors.New("no changes added to commit")
	}
	return false, nil
}

func writeCommit(ctx *context.Context, msg string) error {
	w, err := ctx.Repo.Worktree()
	if err != nil {
		return err
	}

	if _, err := w.Commit(msg, &git.CommitOptions{}); err != nil {
		return err
	}

	return nil
}

func promptEditor(msg string) string {
	prompt := &survey.Editor{
		Message:       "Edit commit message",
		Default:       msg,
		AppendDefault: true,
		FileName:      "COMMIT_EDITMSG*",
	}
	survey.AskOne(prompt, &msg)
	return msg
}
