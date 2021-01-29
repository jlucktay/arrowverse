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

// Package serve holds the logic for the 'arrowverse serve' subcommand.
package serve

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmd returns the 'serve' subcommand.
func NewCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"served", "server", "service", "svc"},

		Short: "Serve scraped data through a web front end",
		Long: `Scrapes data from a wiki website to populate an in-memory collection and then
renders and serves this data through a web front end.`,

		Args: cobra.MaximumNArgs(0),

		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Printf("%s called\n", cmd.CalledAs())
		},
	}

	return serveCmd
}
