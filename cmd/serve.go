package cmd

import (
	"dailies/pkg/server"
	"log"
	"github.com/spf13/cobra"
)

var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Spins up the local REST server backend dashboard mapping engines",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Start(port); err != nil {
			log.Fatalf("Server unexpectedly crashed: %v", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Target running port configuration assignment link")
	RootCmd.AddCommand(serveCmd)
}
