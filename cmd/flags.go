package cmd

import (
	"fmt"
	"log/slog"
	"os"

	dwca "github.com/gnames/dwca/pkg"
	"github.com/gnames/dwca/pkg/config"
	"github.com/spf13/cobra"
)

type flagFunc func(cmd *cobra.Command)

func debugFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("debug")
	if b {
		lopts := &slog.HandlerOptions{Level: slog.LevelDebug}
		handle := slog.NewJSONHandler(os.Stderr, lopts)
		logger := slog.New(handle)
		slog.SetDefault(logger)
	}
}

func rootDirFlag(cmd *cobra.Command) {
	root, _ := cmd.Flags().GetString("root-dir")
	if root != "" {
		opts = append(opts, config.OptRootPath(root))
	}
}

func archiveFlag(cmd *cobra.Command) {
	archive, _ := cmd.Flags().GetString("archive-format")
	if archive != "" {
		opts = append(opts, config.OptArchiveCompression(archive))
	}
}

func csvFlag(cmd *cobra.Command) {
	csv, _ := cmd.Flags().GetString("csv-type")
	if csv != "" {
		opts = append(opts, config.OptOutputCSVType(csv))
	}
}

func jobsNumFlag(cmd *cobra.Command) {
	jobs, _ := cmd.Flags().GetInt("jobs-number")
	if jobs > 0 {
		opts = append(opts, config.OptJobsNum(jobs))
	}
}

func versionFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		version := dwca.Version()
		fmt.Printf(
			"\nVersion: %s\nBuild:   %s\n\n",
			version.Version,
			version.Build,
		)
		os.Exit(0)
	}
}
