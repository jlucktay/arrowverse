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
	"strings"

	"github.com/spf13/cobra"

	"go.jlucktay.dev/arrowverse/pkg/cmd/root/api"
	"go.jlucktay.dev/arrowverse/pkg/cmd/root/scrape"
	"go.jlucktay.dev/arrowverse/pkg/cmd/root/serve"
)

var cfgFile string

func NewCmd() *cobra.Command {
	// cmd represents the base command when called without any subcommands.
	cmd := &cobra.Command{
		Use:   "arrowverse",
		Short: "Arrowverse episode watch order web app",
		Long: `Arrowverse runs a web app to guide visitors and help them keep track with an
episode watch order of all TV shows in DC's Arrowverse.

If no value is specified for the --config flag, the following locations will be
searched in descending order of preference for a file named '` + cfgName + `.` + cfgType + `':
 - the current working directory
 - $HOME/.config/
 - /etc/arrowverse/`,

		Example: "  arrowverse scrape\n" +
			"  arrowverse serve --config=/opt/arrowverse/config." + cfgType,

		Args: cobra.MaximumNArgs(0),

		PersistentPreRunE: initConfig,
		RunE:              checkConfig,

		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 3, //nolint:gomnd // Increase from default of 2
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to "+strings.ToUpper(cfgType)+" configuration file")

	cmd.AddCommand(api.NewCmd())
	cmd.AddCommand(scrape.NewCmd())
	cmd.AddCommand(serve.NewCmd())

	return cmd
}
