// Package cmd is the root command for this app.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.jlucktay.dev/arrowverse/pkg/collection"
	"go.jlucktay.dev/arrowverse/pkg/collection/inmemory"
	"go.jlucktay.dev/arrowverse/pkg/scrape"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arrowverse",
	Short: "Arrowverse episode watch order web app",
	Long: `Arrowverse runs a web app to guide visitors and help them keep track with an
episode watch order of all TV shows in DC's Arrowverse.`,

	Example: "arrowverse --config=/opt/arrowverse/arrowverse-config.json",

	Args: func(_ *cobra.Command, args []string) error {
		if len(args) > 0 {
			sb := strings.Builder{}

			fmt.Fprintln(&sb, "unknown CLI argument(s):")

			for i := range args {
				fmt.Fprintf(&sb, "- %s\n", args[i])
			}

			return errors.New(sb.String())
		}

		return nil
	},

	DisableSuggestions: false,
	SilenceErrors:      false,
	SilenceUsage:       false,

	SuggestionsMinimumDistance: 3,

	Run: func(cmd *cobra.Command, _ []string) {
		episodeLists, errEL := scrape.EpisodeLists()
		if errEL != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "could not get episode lists: %v\n", errEL)
		}

		var cs collection.Shows = &inmemory.Collection{}

		for s, el := range episodeLists {
			show, errEps := scrape.Episodes(s, el)
			if errEps != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "could not get episode details for '%s': %v\n", s, errEps)
			}

			if errAdd := cs.Add(show); errAdd != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "could not add '%s' details to collection: %v\n", show.Name, errAdd)
			}
		}

		includedShows := []string{
			"Arrow",
			"Batwoman",
			"Black Lightning",
			"Constantine",
			"DC's Legends of Tomorrow",
			"Freedom Fighters: The Ray",
			"Supergirl",
			"The Flash (The CW)",
			"Vixen",
		}

		// Excluded shows:
		// "Birds of Prey"
		// "The Flash (CBS)"

		csio, errIO := cs.InOrder(includedShows...)
		if errIO != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "could not get episode details: %v\n", errIO)
		}

		for i := range csio {
			fmt.Fprintf(cmd.OutOrStdout(), "[%03d] %s\n", i, csio[i])
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.arrowverse.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("arrowverse")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".arrowverse" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".arrowverse")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
