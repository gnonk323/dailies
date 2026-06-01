package cmd

import (
  "fmt"
  "os"
  "github.com/spf13/cobra"
  "dailies/pkg/sync" 
)

var pullCmd = &cobra.Command{
  Use:   "pull",
  Short: "Pull the latest daily logs from the remote GitHub repository",
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Checking for remote updates...")

    if err := sync.SyncPull(); err != nil {
      fmt.Printf("Error during pull sync: %v\n", err)
      os.Exit(1)
    }

    fmt.Println("Data sync complete.")
  },
}

func init() {
  RootCmd.AddCommand(pullCmd)
}
