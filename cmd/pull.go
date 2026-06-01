package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull the latest daily logs from the remote GitHub repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fetching latest changes from remote...")

		gitCmd := exec.Command("git", "pull", "origin", "main")
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		if err := gitCmd.Run(); err != nil {
			fmt.Printf("Error pulling data: %v\n", err)
			return
		}

		fmt.Println("Data successfully synced from remote.")
	},
}

func init() {
	RootCmd.AddCommand(pullCmd)
}
