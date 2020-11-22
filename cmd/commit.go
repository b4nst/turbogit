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
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/b4nst/turbogit/internal/format"
	intgit "github.com/b4nst/turbogit/internal/git"
	git "github.com/libgit2/git2go/v30"
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

type CommitCmdOption struct {
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
	// Current repository
	Repo *git.Repository
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
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	ValidArgs:    format.AllCommitType(),
	Run:          runCommitCmd,
}

func runCommitCmd(cmd *cobra.Command, args []string) {
	cco, err := parseCommitCmd(cmd, args)
	if err != nil {
		log.Fatal(err)
	}
	err = runCommit(cco)
	if err != nil {
		log.Fatal(err)
	}
}

func parseCommitCmd(cmd *cobra.Command, args []string) (*CommitCmdOption, error) {
	// --type
	fType, err := cmd.Flags().GetString("type")
	if err != nil {
		return nil, err
	}
	if fType == "" {
		fType = args[0]
		args = args[1:]
	}
	ctype := format.FindCommitType(fType)

	// --breaking-changes
	fBreakingChanges, err := cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return nil, err
	}

	// --scope
	fScope, err := cmd.Flags().GetString("scope")
	if err != nil {
		return nil, err
	}

	// --edit
	fEdit, err := cmd.Flags().GetBool("edit")
	if err != nil {
		return nil, err
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

	return &CommitCmdOption{
		CType:           ctype,
		BreakingChanges: fBreakingChanges,
		Message:         strings.Join(args, " "),
		PromptEditor:    fEdit,
		Scope:           fScope,
		Repo:            repo,
	}, nil
}

func runCommit(cco *CommitCmdOption) error {
	// Check if working tree is clean
	nc, err := needCommit(cco.Repo)
	if err != nil {
		return err
	}
	if !nc {
		fmt.Println("Nothing to commit, working tree clean")
		return nil
	}

	err = intgit.PreCommitHook(cco.Repo.Path())
	if err != nil {
		return fmt.Errorf("Error during pre-commit hook: %s", err.Error())
	}

	msg, err := intgit.PrepareCommitMsgHook(cco.Repo.Path())
	if err != nil {
		return fmt.Errorf("Error during prepare-commit-msg hook: %s", err.Error())
	}
	if msg == "" {
		msg = cco.Message
	}
	cmo := format.ParseCommitMsg(msg)
	if cmo == nil {
		// Parse commit type
		ctype := cco.CType
		if ctype == format.NilCommit {
			return errors.New("A commit type is required")
		}
		cmo = &format.CommitMessageOption{
			Ctype: ctype, BreakingChanges: cco.BreakingChanges, Description: msg, Scope: cco.Scope,
		}
	}

	cmsg := format.CommitMessage(cmo)
	if cco.PromptEditor {
		cmsg = promptEditor(cmsg)
	}
	cmsg, err = intgit.CommitMsgHook(cco.Repo.Path(), cmsg)
	if err != nil {
		return fmt.Errorf("Error during commit-msg hook: %s", err.Error())
	}

	// // Write commit
	commit, err := writeCommit(cco.Repo, cmsg)
	if err != nil {
		return err
	}
	h, err := commit.ShortId()
	if err != nil {
		return err
	}
	fmt.Println(h, commit.Summary())

	err = intgit.PostCommitHook(cco.Repo.Path())
	if err != nil {
		fmt.Println("Warning, post-commit hook failed:", err.Error())
	}

	return nil
}

func needCommit(r *git.Repository) (bool, error) {
	s, err := r.StatusList(&git.StatusOptions{Show: git.StatusShowIndexAndWorkdir, Flags: git.StatusOptIncludeUntracked})
	if err != nil {
		return false, err
	}

	count, err := s.EntryCount()
	if err != nil {
		return false, err
	}
	if count <= 0 {
		return false, nil
	}
	for i := 0; i < count; i++ {
		se, err := s.ByIndex(i)
		if err != nil {
			return false, err
		}
		if se.Status <= git.StatusIndexTypeChange {
			return true, nil
		}
	}
	return false, errors.New("No changes added to commit")
}

func writeCommit(r *git.Repository, msg string) (*git.Commit, error) {
	sig, err := signature(r)
	if err != nil {
		return nil, err
	}

	idx, err := r.Index()
	if err != nil {
		return nil, err
	}
	treeId, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}
	tree, err := r.LookupTree(treeId)
	if err != nil {
		return nil, err
	}

	parents := []*git.Commit{}
	head, err := r.Head()
	if err == nil { // We found head
		headRef, err := r.LookupCommit(head.Target())
		if err != nil {
			return nil, err
		}
		parents = append(parents, headRef)
	}

	oid, err := r.CreateCommit("HEAD", sig, sig, msg, tree, parents...)
	if err != nil {
		return nil, err
	}

	return r.LookupCommit(oid)
}

func signature(r *git.Repository) (*git.Signature, error) {
	config, err := r.Config()
	if err != nil {
		return nil, err
	}

	email, _ := config.LookupString("user.email")
	name, _ := config.LookupString("user.name")

	return &git.Signature{
		Email: email,
		Name:  name,
		When:  time.Now(),
	}, nil
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
