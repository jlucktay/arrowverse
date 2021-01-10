// Package cmd is the root command for this app.
package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
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

	RunE: runRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
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
