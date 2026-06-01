package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Stage, commit, and push only the data directory to GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Syncing local logs to GitHub...")

		if err := runGit("add", "data/"); err != nil {
			fmt.Printf("Error staging data: %v\n", err)
			return
		}

		statusCmd := exec.Command("git", "diff", "--cached", "--quiet")
		if err := statusCmd.Run(); err == nil {
			fmt.Println("No new data changes to push.")
			return
		}

		commitMsg := fmt.Sprintf("site: sync daily logs - %s", time.Now().Format("2006-01-02 15:04"))
		if err := runGit("commit", "-m", commitMsg); err != nil {
			fmt.Printf("Error committing data: %v\n", err)
			return
		}

		if err := runGit("push", "origin", "main"); err != nil {
			fmt.Printf("Error pushing data: %v\n", err)
			return
		}

		fmt.Println("Data successfully backed up to GitHub.")
	},
}

func runGit(args ...string) error {
	gitCmd := exec.Command("git", args...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	return gitCmd.Run()
}

func init() {
	RootCmd.AddCommand(pushCmd)
}
