package cmdbuilder

import "github.com/spf13/cobra"

// AppendPreRun appends a prerun function to cmd.
func AppendPreRun(cmd *cobra.Command, prerun func(cmd *cobra.Command, args []string)) {
	mem := cmd.PreRun

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if mem != nil {
			mem(cmd, args)
		}

		prerun(cmd, args)
	}
}
