package sync

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// SyncPull checks if we are behind remote and fast-forwards changes
func SyncPull() error {
	// quietly fetch remote updates
	if err := exec.Command("git", "fetch", "origin").Run(); err != nil {
		return errors.New("failed to fetch from GitHub; check connection")
	}

	// check if local master is behind origin/master
	cmd := exec.Command("git", "rev-list", "HEAD..origin/master", "--count")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		countStr := out.String()
		// if there are incoming commits, fast-forward pull
		if countStr != "0\n" && countStr != "" {
			fmt.Println("Local data behind remote. Syncing down latest changes...")
			pullCmd := exec.Command("git", "pull", "--ff-only", "origin", "master")
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return errors.New("Conflict detected! Local logs diverge from GitHub; resolve manually")
			}
		}
	}
	return nil
}

// SyncPush stages, commits, and pushes ONLY the data/ directory.
func SyncPush() error {
	fmt.Println("Backing up data to GitHub...")

	// stage only the data folder
	if err := exec.Command("git", "add", "data/").Run(); err != nil {
		return fmt.Errorf("failed to stage data folder: %w", err)
	}

	// check if there are actual changes staged to avoid empty commit errors
	if err := exec.Command("git", "diff", "--cached", "--quiet").Run(); err == nil {
		fmt.Println("No new changes to push.")
		return nil
	}

	// commit changes
	commitMsg := fmt.Sprintf("site: auto-sync daily logs - %s", time.Now().Format("2006-01-02 15:04"))
	if err := exec.Command("git", "commit", "-m", commitMsg).Run(); err != nil {
		return fmt.Errorf("failed to commit logs: %w", err)
	}

	// push to remote
	pushCmd := exec.Command("git", "push", "origin", "master")
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("Failed to push to GitHub: %w", err)
	}

	fmt.Println("Data successfully pushed.")
	return nil
}
