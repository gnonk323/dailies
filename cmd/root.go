package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var targetDate string

var RootCmd = &cobra.Command{
	Use:   "dailies",
	Short: "A local-first personal daily tracker CLI",
}

func Execute() {
	RootCmd.PersistentFlags().StringVarP(&targetDate, "date", "d", "", "Specify target date override (YYYY-MM-DD)")

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
