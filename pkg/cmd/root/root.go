/*
Copyright © 2021 James Lucktaylor <jlucktay@gmail.com>

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

// Package root holds the logic for the root 'arrowverse' command.
package root

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"go.jlucktay.dev/arrowverse/pkg/cmd/root/api"
	"go.jlucktay.dev/arrowverse/pkg/cmd/root/scrape"
	"go.jlucktay.dev/arrowverse/pkg/cmd/root/serve"
)

func NewCmd() *cobra.Command {
	// cmd represents the base command when called without any subcommands.
	var cmd = &cobra.Command{
		Use:   "arrowverse",
		Short: "Arrowverse episode watch order web app",
		Long: `Arrowverse runs a web app to guide visitors and help them keep track with an
	episode watch order of all TV shows in DC's Arrowverse.`,

		Example: "arrowverse --config=/opt/arrowverse/arrowverse-config.json",

		Args: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 {
				sb := strings.Builder{}

				for i := range args {
					fmt.Fprintf(&sb, "- %s\n", args[i])
				}

				return fmt.Errorf("%w:\n%s", ErrUnknownCLIArguments, sb.String())
			}

			return nil
		},

		DisableSuggestions: false,
		SilenceErrors:      false,
		SilenceUsage:       false,

		SuggestionsMinimumDistance: 3, //nolint:gomnd // Increase from default of 2
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.arrowverse.yaml)")

	cobra.OnInitialize(initConfig)

	cmd.AddCommand(api.NewCmd())
	cmd.AddCommand(scrape.NewCmd())
	cmd.AddCommand(serve.NewCmd())

	return cmd
}
