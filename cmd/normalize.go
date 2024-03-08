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
	"log/slog"
	"os"

	dwca "github.com/gnames/dwca/pkg"
	"github.com/gnames/dwca/pkg/config"
	"github.com/spf13/cobra"
)

// normalizeCmd represents the normalize command
var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalizes some known ambiguities in DwCA files.",
	Long: `There are some known ambiguities in DwCA files, that are
	normalized by this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := []flagFunc{
			debugFlag, rootDirFlag, jobsNumFlag, archiveFlag, csvFlag,
		}
		for _, v := range flags {
			v(cmd)
		}
		in, out := getInput(cmd, args)
		cfg := config.New(opts...)
		arc, err := dwca.Factory(in, cfg)
		if err != nil {
			slog.Error("Cannot initialize DwCA", "error", err)
			os.Exit(1)
		}

		slog.Info(
			"Configuration",
			"concurrent_jobs", cfg.JobsNum,
			"output_csv_type", cfg.OutputCSVType,
			"archive_type", cfg.OutputArchiveCompression)

		err = arc.Load(cfg.ExtractPath)
		if err != nil {
			slog.Error("Cannot load DwCA", "error", err)
			os.Exit(1)
		}

		err = arc.Normalize()
		if err != nil {
			slog.Error("Cannot normalize DwCA", "error", err)
			os.Exit(1)
		}

		if arc.Config().OutputArchiveCompression == "zip" {
			out += ".zip"
			err = arc.ZipNormalized(out)
		} else {
			out += ".tar.gz"
			err = arc.TarGzNormalized(out)
		}
		if err != nil {
			slog.Error("Cannot archive DwCA data", "error", err)
			os.Exit(1)
		}

		slog.Info("DwCA normalized", "input", in, "output", out)
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// normalizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	normalizeCmd.Flags().StringP("archive-format", "a", "",
		"format of the output archive (tar or zip)",
	)

	normalizeCmd.Flags().StringP("csv-type", "c", "",
		"type of CSV files in the output archive (csv or tsv)",
	)
}

func getInput(cmd *cobra.Command, args []string) (in, out string) {
	switch len(args) {
	case 1:
		path := args[0]
		return path, path + ".norm"
	case 2:
		return args[0], args[1]
	default:
		_ = cmd.Help()
		os.Exit(0)
	}
	return "", ""
}
