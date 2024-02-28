/*
Copyright © 2024 Dmitry Mozzherin <dmozzherin@gmail.com>

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
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/dwca/pkg/config"
	"github.com/gnames/gnsys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed dwca.yaml
var configText string

type configData struct {
	RootPath                 string
	OutputArchiveCompression string
	OutputCSVType            string
	JobsNum                  int
}

var opts []config.Option

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dwca",
	Short: "Utities for Darwin Core Archives (DwC-A) processing.",
	Long: `dwca is a command line tool for processing Darwin Core Archives
	(DwC-A).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, _ []string) {
		versionFlag(cmd)
		flags := []flagFunc{debugFlag}
		for _, v := range flags {
			v(cmd)
		}
		_ = cmd.Help()
	},
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().IntP(
		"jobs-number", "j", 0,
		"number of concurrent jobs",
	)

	rootCmd.PersistentFlags().StringP(
		"root-dir", "r", "",
		"root path for the DwCA file",
	)

	rootCmd.PersistentFlags().BoolP(
		"debug", "d", false,
		"debug mode",
	)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "V", false, "show version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var configDir string
	var err error
	configFile := "dwca"
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	configDir = filepath.Join(home, ".config")

	// Search config in home directory with name ".gnmatcher" (without extension).
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)

	_ = viper.BindEnv("RootPath", "DWCA_ROOT_PATH")
	_ = viper.BindEnv("OutputArchiveCompression", "DWCA_OUTPUT_ARCHIVE_COMPRESSION")
	_ = viper.BindEnv("OutputCSVType", "DWCA_OUTPUT_CSV_TYPE")
	_ = viper.BindEnv("JobsNum", "DWCA_JOBS_NUM")

	viper.AutomaticEnv() // read in environment variables that match

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yaml", configFile))
	touchConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		msg := fmt.Sprintf("Using config file %s.", viper.ConfigFileUsed())
		slog.Info(msg)
	}

	getOpts()
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts() {
	cfgCli := &configData{}
	err := viper.Unmarshal(cfgCli)
	if err != nil {
		msg := fmt.Sprintf("Cannot deserialize config file: %s", err)
		slog.Error(msg)
		os.Exit(1)
	}

	if cfgCli.RootPath != "" {
		opts = append(opts, config.OptRootPath(cfgCli.RootPath))
	}

	if cfgCli.OutputArchiveCompression != "" {
		opts = append(
			opts,
			config.OptArchiveCompression(cfgCli.OutputArchiveCompression),
		)
	}

	if cfgCli.OutputCSVType != "" {
		opts = append(
			opts,
			config.OptOutputCSVType(cfgCli.OutputCSVType),
		)
	}

	if cfgCli.JobsNum != 0 {
		opts = append(opts, config.OptJobsNum(cfgCli.JobsNum))
	}
}

// touchConfigFile checks if config file exists, and if not, it gets created.
func touchConfigFile(configPath string) error {
	exists, err := gnsys.FileExists(configPath)
	if exists || err != nil {
		return err
	}

	msg := fmt.Sprintf("Creating config file '%s'", configPath)
	slog.Info(msg)
	createConfig(configPath)
	return nil
}

// createConfig creates config file.
func createConfig(path string) {
	err := gnsys.MakeDir(filepath.Dir(path))
	if err != nil {
		msg := fmt.Sprintf("Cannot create dir %s: %s", path, err)
		slog.Error(msg)
		os.Exit(1)
	}

	err = os.WriteFile(path, []byte(configText), 0644)
	if err != nil {
		msg := fmt.Sprintf("Cannot write to file %s: %s", path, err)
		slog.Error(msg)
		os.Exit(1)
	}
}
