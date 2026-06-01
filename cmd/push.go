package cmd

import (
  "fmt"
  "os"
  "github.com/spf13/cobra"
  "dailies/pkg/sync" 
)

var pushCmd = &cobra.Command{
  Use:   "push",
  Short: "Stage, commit, and push only the data directory to GitHub",
  Run: func(cmd *cobra.Command, args []string) {
    if err := sync.SyncPush(); err != nil {
      fmt.Printf("Error during push sync: %v\n", err)
      os.Exit(1)
    }
  },
}

func init() {
  RootCmd.AddCommand(pushCmd)
}
