/*
Copyright © 2022 Robert Sigler <sigler@improvisedscience.org>

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
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openweather",
	Short: "Fetch weather data from OpenWeatherMap",
	Long: `Fetch weather data from OpenWeatherMap. 
An API key is required. See https://openweathermap.org/ for more information. 
API Key, lat, and lon can be provided via flags or via a config file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	var err error

	// Find home directory.
	homeConfigDir, err = os.UserConfigDir()
	cobra.CheckErr(err)
	homeConfigDir = path.Join(homeConfigDir, "openweather")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s/config.yaml)", homeConfigDir))
	rootCmd.Flags().String("apikey", "", "OpenWeatherMap API key")
	rootCmd.Flags().Float64("lat", 0.0, "latitude")
	rootCmd.Flags().Float64("lon", 0.0, "longitude")

	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("lat", rootCmd.PersistentFlags().Lookup("lat"))
	viper.BindPFlag("lon", rootCmd.PersistentFlags().Lookup("lon"))

	rootCmd.MarkFlagRequired("apikey")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".openweather" (without extension).
		viper.AddConfigPath(homeConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
