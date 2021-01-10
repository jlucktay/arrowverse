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

package root

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgName = "arrowverse"
	cfgType = "json"
)

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command, _ []string) error {
	viper.SetEnvPrefix(cfgName)
	viper.SetConfigType(cfgType) // REQUIRED if the config file does not have the extension in the name

	if cfgFile != "" {
		// If a config file is specified explicitly with the flag, use that.
		viper.SetConfigFile(cfgFile)
	} else {
		// Otherwise, look in a few different places.

		// Find home directory.
		home, errHomeDir := homedir.Dir()
		if errHomeDir != nil {
			return errHomeDir
		}

		// Expect config file name to be "arrowverse.<something>".
		viper.SetConfigName(cfgName)

		// Search for a config file in the following order of preference (most preferred first):
		// - the current working directory
		viper.AddConfigPath(".")
		// - user-specific config directory per the XDG Base Directory Specification
		// cf. https://specifications.freedesktop.org/basedir-spec/latest/ar01s03.html
		viper.AddConfigPath(home + "/.config/")
		// - '/etc/arrowverse/'
		viper.AddConfigPath("/etc/" + cfgName + "/")
	}


	// Read in environment variables that match 'ARROWVERSE_*'.
	viper.AutomaticEnv()

	// If a config file is found, (attempt to) read it in.
	if errViperRead := viper.ReadInConfig(); errViperRead != nil {
		_, osPE := errViperRead.(*os.PathError)
		_, viperCFNFE := errViperRead.(viper.ConfigFileNotFoundError)

		if osPE || viperCFNFE {
			if cfgFile != "" {
				return fmt.Errorf("Config file '%s' could not be read at specified path: %w", cfgFile, errViperRead)
			}

			return nil
		}

		return fmt.Errorf("Config file '%s' was found but another error occurred:\n%w",
			viper.ConfigFileUsed(), errViperRead)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Config file found and successfully parsed: %s\n", viper.ConfigFileUsed())

	// Read an update to a config file while running and not miss a beat.
	viper.WatchConfig()

	return nil
}

// checkConfig goes through configured settings and prints them out.
func checkConfig(cmd *cobra.Command, _ []string) error {
	vas := viper.AllSettings()

	if len(vas) == 0 {
		return errors.New("There are no established configuration values.")
	}

	fmt.Fprintln(cmd.OutOrStdout(), "The following configuration values were established:")

	for k, v := range viper.AllSettings() {
		fmt.Fprintf(cmd.OutOrStdout(), "key: '%s', value (of type %[2]T): %#[2]v\n", k, v)
	}

	return nil
}
