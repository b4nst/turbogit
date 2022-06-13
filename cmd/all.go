package cmd

import (
	cc "github.com/b4nst/turbogit/cmd/git-cc/cmd"
	check "github.com/b4nst/turbogit/cmd/git-check/cmd"
	logs "github.com/b4nst/turbogit/cmd/git-logs/cmd"
	gnew "github.com/b4nst/turbogit/cmd/git-new/cmd"
	release "github.com/b4nst/turbogit/cmd/git-release/cmd"
	"github.com/spf13/cobra"
)

var Commands = []*cobra.Command{
	cc.RootCmd,
	check.RootCmd,
	logs.RootCmd,
	gnew.RootCmd,
	release.RootCmd,
}

var Command = &cobra.Command{
	Use:   "git",
	Short: "Turbogit is a set of opinionated git subcommand",
}

func init() {
	Command.AddCommand(Commands...)
}
