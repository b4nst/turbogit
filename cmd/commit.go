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
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	tugit "github.com/b4nst/turbogit/pkg/git"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(CommitCmd)

	cmdbuilder.RepoAware(CommitCmd)

	CommitCmd.Flags().StringP("type", "t", "", fmt.Sprintf("Commit types %s", format.AllCommitType()))
	CommitCmd.RegisterFlagCompletionFunc("type", typeFlagCompletion)
	CommitCmd.Flags().BoolP("breaking-changes", "c", false, "Commit contains breaking changes")
	CommitCmd.Flags().BoolP("edit", "e", false, "Prompt editor to edit your message (add body or/and footer(s))")
	CommitCmd.Flags().StringP("scope", "s", "", "Add a scope")
	CommitCmd.Flags().BoolP("amend", "a", false, "Amend commit")
}

func typeFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return format.AllCommitType(), cobra.ShellCompDirectiveDefault
}

var CommitCmd = &cobra.Command{
	Use:   "commit [type] [subject]",
	Short: "Commit using conventional commit message",
	// DisableFlagsInUseLine: true,
	Example: `
# Commit a new feature (feat: a new feature)
$ tug commit feat a new feature

# Commit a fix that brings breaking changes (fix!: API break)
$ tug commit fix -c API break

# Add a scope to the commit (refactor(scope): a scopped refactor)
$ tug commit refactor a scopped refactor -s scope

# Open your editor to edit the commit message
$ tug commit ci -e message

# Ammend last commit type
$ tug commit -a -t fix
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		amend, _ := cmd.Flags().GetBool("amend")
		fType, _ := cmd.Flags().GetString("type")
		ma := 0
		if !amend {
			ma++ // Need at least one argument for the description
			if fType == "" {
				ma++ // Also need a type since it was not passed as flag
			}
		}
		if len(args) < ma {
			return fmt.Errorf("requires at least %d arg(s), only received %d", ma, len(args))
		}
		return nil
	},
	// SilenceUsage: true,
	ValidArgs: format.AllCommitType(),

	Run: func(cmd *cobra.Command, args []string) {
		cco, err := parseCommitCmd(cmd, args)
		cobra.CheckErr(err)
		cobra.CheckErr(runCommit(cco))
	},
}

type commitOpt struct {
	// Commit type
	CType format.CommitType
	// True if this commit introduces breaking changes
	BreakingChanges bool
	// Should prompt an editor before committing
	PromptEditor bool
	// Commit scope (optional)
	Scope string
	// Commit message
	Message string
	// Amend
	Amend bool
	// Current repository
	Repo *git.Repository
}

func parseCommitCmd(cmd *cobra.Command, args []string) (*commitOpt, error) {
	opt := &commitOpt{}
	var err error

	// --type
	fType, err := cmd.Flags().GetString("type")
	if err != nil {
		return nil, err
	}
	opt.CType = format.FindCommitType(fType)
	if opt.CType == format.NilCommit && len(args) > 0 {
		opt.CType = format.FindCommitType(args[0])
		if opt.CType != format.NilCommit {
			// Type was in first arg
			args = args[1:]
		}
	}

	// --breaking-changes
	opt.BreakingChanges, err = cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return nil, err
	}

	// --scope
	opt.Scope, err = cmd.Flags().GetString("scope")
	if err != nil {
		return nil, err
	}

	// --edit
	opt.PromptEditor, err = cmd.Flags().GetBool("edit")
	if err != nil {
		return nil, err
	}

	// --amend
	opt.Amend, err = cmd.Flags().GetBool("amend")
	if err != nil {
		return nil, err
	}

	// Find repo
	opt.Repo = cmdbuilder.GetRepo(cmd)

	opt.Message = strings.Join(args, " ")

	return opt, nil
}

func runCommit(cco *commitOpt) error {
	var ca *git.Commit
	if cco.Amend {
		o, err := cco.Repo.RevparseSingle("HEAD")
		if err != nil {
			return err
		}
		ca, err = o.AsCommit()
		if err != nil {
			return err
		}
	}

	// Check if working tree is clean
	if ca == nil { // Only when not amending
		nc, err := tugit.StageReady(cco.Repo)
		if err != nil {
			return err
		}
		if !nc {
			fmt.Println("Nothing to commit, working tree clean")
			return nil
		}
	}

	msg := ""
	if ca != nil {
		// Init message from amend origin
		msg = ca.Message()
	} else {
		// Init message from hooks
		if err := tugit.PreCommitHook(cco.Repo.Path()); err != nil {
			return fmt.Errorf("Error during pre-commit hook: %s", err.Error())
		}
		m, err := tugit.PrepareCommitMsgHook(cco.Repo.Path())
		if err != nil {
			return fmt.Errorf("Error during prepare-commit-msg hook: %s", err.Error())
		}
		msg = m
	}
	// Try to parse the initial message
	cmo := format.ParseCommitMsg(msg)
	if cmo == nil {
		// If not formatted put raw message as Description
		cmo = &format.CommitMessageOption{Description: msg}
	}
	if err := cmo.Overwrite(&format.CommitMessageOption{
		Ctype:           cco.CType,
		BreakingChanges: cco.BreakingChanges,
		Description:     cco.Message,
		Scope:           cco.Scope,
	}); err != nil {
		return err
	}
	if err := cmo.Check(); err != nil {
		return err
	}

	// Build commit message
	cmsg := format.CommitMessage(cmo)
	if cco.PromptEditor {
		cmsg = promptEditor(cmsg)
	}
	var err error
	cmsg, err = tugit.CommitMsgHook(cco.Repo.Path(), cmsg)
	if err != nil {
		return fmt.Errorf("Error during commit-msg hook: %s", err.Error())
	}

	// Write commit
	var commit *git.Commit
	if ca != nil {
		commit, err = tugit.Amend(ca, cmsg)
	} else {
		commit, err = tugit.Commit(cco.Repo, cmsg)
	}
	if err != nil {
		return err
	}

	h, err := commit.ShortId()
	if err != nil {
		return err
	}
	fmt.Println(h, commit.Summary())

	err = tugit.PostCommitHook(cco.Repo.Path())
	if err != nil {
		fmt.Println("Warning, post-commit hook failed:", err.Error())
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
