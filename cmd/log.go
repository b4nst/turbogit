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
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/araddon/dateparse"
	"github.com/b4nst/turbogit/internal/cmdbuilder"
	"github.com/b4nst/turbogit/pkg/format"
	git "github.com/libgit2/git2go/v33"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(LogCmd)

	cmdbuilder.RepoAware(LogCmd)

	LogCmd.Flags().BoolP("all", "a", false, "Pretend as if all the refs in refs/, along with HEAD, are listed on the command line as <commit>. If set on true, the --from option will be ignored.")
	LogCmd.Flags().Bool("no-color", false, "Disable color output")
	LogCmd.Flags().StringP("from", "f", "HEAD", "Logs only commits reachable from this one")
	LogCmd.Flags().String("since", "", "Show commits more recent than a specific date")
	LogCmd.Flags().String("until", "", "Show commits older than a specific date")
	// logCmd.Flags().String("path", "", "Filter commits based on the path of files that are updated. Accept regexp")
	// Filters
	LogCmd.Flags().StringArrayP("type", "t", []string{}, "Filter commits by type (repeatable option)")
	LogCmd.Flags().StringArrayP("scope", "s", []string{}, "Filter commits by scope (repeatable option)")
	LogCmd.Flags().BoolP("breaking-changes", "c", false, "Only shows breaking changes")
}

// LogCmd represents the log command
var LogCmd = &cobra.Command{
	Use:   "logs",
	Short: "Shows the commit logs.",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		opt, err := parseCmd(cmd, args)
		cobra.CheckErr(err)
		cobra.CheckErr(runLog(opt))
	},
}

type logOpt struct {
	All            bool
	NoColor        bool
	From           string
	Since          *time.Time
	Until          *time.Time
	Types          []format.CommitType
	Scopes         []string
	BreakingChange bool
	Repo           *git.Repository
}

func parseCmd(cmd *cobra.Command, args []string) (*logOpt, error) {
	opt := &logOpt{}
	var err error

	// --all
	opt.All, err = cmd.Flags().GetBool("all")
	if err != nil {
		return nil, err
	}
	// --no-color
	opt.NoColor, err = cmd.Flags().GetBool("no-color")
	if err != nil {
		return nil, err
	}
	// --from
	opt.From, err = cmd.Flags().GetString("from")
	if err != nil {
		return nil, err
	}
	// --since
	fSince, err := cmd.Flags().GetString("since")
	if err != nil {
		return nil, err
	}
	if fSince != "" {
		date, err := dateparse.ParseAny(fSince)
		if err != nil {
			return nil, err
		}
		opt.Since = &date
	}
	// --until
	fUntil, err := cmd.Flags().GetString("until")
	if err != nil {
		return nil, err
	}
	if fUntil != "" {
		date, err := dateparse.ParseAny(fUntil)
		if err != nil {
			return nil, err
		}
		opt.Until = &date
	}
	// --types
	fTypes, err := cmd.Flags().GetStringArray("type")
	if err != nil {
		return nil, err

	}
	for _, v := range fTypes {
		opt.Types = append(opt.Types, format.FindCommitType(v))
		// TODO warn or error on nil commit type
	}
	// --scopes
	opt.Scopes, err = cmd.Flags().GetStringArray("scope")
	if err != nil {
		return nil, err
	}
	// --breaking-changes
	opt.BreakingChange, err = cmd.Flags().GetBool("breaking-changes")
	if err != nil {
		return nil, err
	}

	opt.Repo = cmdbuilder.GetRepo(cmd)

	return opt, nil
}

func runLog(opt *logOpt) error {
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

	// Build filters
	filters := []LogFilter{
		Since(opt.Since),
		Until(opt.Until),
		Type(opt.Types),
		Scope(opt.Scopes),
		BreakingChange(opt.BreakingChange),
	}

	tw := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', 0)
	defer tw.Flush()
	if err := walk.Iterate(buildLogWalker(tw, !opt.NoColor, filters)); err != nil {
		return err
	}

	return nil
}

func buildLogWalker(w io.Writer, color bool, filters []LogFilter) func(c *git.Commit) bool {
	return func(c *git.Commit) bool {
		co := format.ParseCommitMsg(c.Message())
		if co == nil {
			co = &format.CommitMessageOption{}
		}
		keep, walk := ApplyFilters(c, co, filters...)
		if !keep {
			return walk
		}

		// Hash
		h, err := c.ShortId()
		if err != nil {
			h = c.Id().String()
		}
		// type
		var ctype string
		if color {
			h = fmt.Sprintf("\x1b[38;5;231m%s\x1b[0m", h)
			ctype = co.Ctype.ColorString()
		} else {
			ctype = co.Ctype.String()
		}
		// description
		msg := co.Description
		if msg == "" {
			msg = c.Summary()
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", h, ctype, msg)

		return walk
	}
}
