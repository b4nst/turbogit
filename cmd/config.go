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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/b4nst/turbogit/internal/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:           "config [key] [value]",
	Short:         "Read or write config.",
	Long:          `If value is not provided display the current value for this key.`,
	Args:          cobra.RangeArgs(1, 2),
	SilenceUsage:  true,
	SilenceErrors: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		context.FromCommand(cmd)
	},
	RunE: configure,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolP("delete", "d", false, "Delete config.")
}

func configure(cmd *cobra.Command, args []string) error {
	delete, err := cmd.Flags().GetBool("delete")
	if err != nil {
		return err
	}

	if len(args) > 1 || delete {
		// Set config
		// TODO do not use viper.WriteConfig. Better manage delete.
		if delete {
			viper.Set(args[0], "")
		} else {
			v, err := parseValue(args[1])
			if err != nil {
				return err
			}
			viper.Set(args[0], v)
		}
		// Write config
		if err := viper.WriteConfig(); err != nil {
			return err
		}
		fmt.Println("Written to :", viper.ConfigFileUsed())
	} else {
		v := viper.Get(args[0])
		if v == nil {
			return fmt.Errorf("[%s] is not set.", args[0])
		} else {
			fmt.Println(v)
		}
	}
	return nil
}

func parseValue(s string) (interface{}, error) {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		var res []interface{}
		for _, ss := range strings.Split(string(s[1:len(s)-1]), ",") {
			v, err := parseValue(strings.Trim(ss, " "))
			if err != nil {
				return nil, err
			}

			res = append(res, v)
		}
		return res, nil
	}

	if strings.HasPrefix(s, "'") {
		return strings.Trim(s, "'"), nil
	}

	re := regexp.MustCompile(`^\d+$`)
	if re.Match([]byte(s)) {
		return strconv.ParseInt(s, 10, 0)
	}

	re = regexp.MustCompile(`(?i)^true|false$`)
	if re.Match([]byte(s)) {
		return strconv.ParseBool(s)
	}

	re = regexp.MustCompile(`^\d+(\.|,)\d+$`)
	if re.Match([]byte(s)) {
		return strconv.ParseFloat(s, 32)
	}

	return s, nil
}
