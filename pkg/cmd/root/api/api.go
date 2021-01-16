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

// Package api holds the logic for the 'arrowverse api' subcommand.
package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

// Options is a struct to support the 'api' subcommand.
type Options struct {
	Port uint16
}

// NewOptions returns initialised Options with defaults set.
func NewOptions() *Options {
	return &Options{Port: 3000}
}

// NewCmd encapsulates the 'api' subcommand.
func NewCmd() *cobra.Command {
	o := NewOptions()

	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Run up a headless web API serving scraped data",
		Long: `Scrapes data from a wiki website to populate an in-memory collection and then
serves this data through a RESTful web API.`,

		Args: cobra.MaximumNArgs(0),

		RunE: func(_ *cobra.Command, _ []string) error {
			app := fiber.New()

			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("Hello, World 👋!")
			})

			listenAddr := fmt.Sprintf(":%d", o.Port)

			return app.Listen(listenAddr)
		},
	}

	apiCmd.Flags().Uint16VarP(&o.Port, "port", "p", o.Port, "port to listen on")

	return apiCmd
}
