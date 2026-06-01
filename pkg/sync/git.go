package sync

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
	"strconv"
	"dailies/pkg/storage"
)

func isGitEnabled() bool {
	enabled, err := strconv.ParseBool(os.Getenv("AUTO_GIT"))
	if err != nil {
		return false
	}
	return enabled
}

// SyncPull checks if we are behind remote and fast-forwards changes
func SyncPull() error {
	enabled := isGitEnabled()
	if !enabled {
		return nil
	}

	dataDir := storage.GetDataDirectory()

	// quietly fetch remote updates
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = dataDir
	if err := fetchCmd.Run(); err != nil {
		return errors.New("failed to fetch from GitHub; check connection")
	}

	// check if local is behind its upstream tracking branch
	cmd := exec.Command("git", "rev-list", "HEAD..@{u}", "--count")
	cmd.Dir = dataDir
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		countStr := out.String()
		// if there are incoming commits, fast-forward pull
		if countStr != "0\n" && countStr != "" {
			fmt.Println("Local data behind remote. Syncing down latest changes...")
			
			pullCmd := exec.Command("git", "pull", "--ff-only")
			pullCmd.Dir = dataDir
			pullCmd.Stdout = os.Stdout
			pullCmd.Stderr = os.Stderr
			if err := pullCmd.Run(); err != nil {
				return errors.New("Conflict detected! Local logs diverge from GitHub; resolve manually")
			}
		}
	}
	return nil
}

// SyncPush stages, commits, and pushes the centralized data directory contents.
func SyncPush() error {
	enabled := isGitEnabled()
	if !enabled {
		return nil
	}

	fmt.Println("Backing up data to GitHub...")
	dataDir := storage.GetDataDirectory()

	// stage everything in the dedicated data repository
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = dataDir
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage data changes: %w", err)
	}

	// check if there are actual changes staged to avoid empty commit errors
	diffCmd := exec.Command("git", "diff", "--cached", "--quiet")
	diffCmd.Dir = dataDir
	if err := diffCmd.Run(); err == nil {
		fmt.Println("No new changes to push.")
		return nil
	}

	// commit changes
	commitMsg := fmt.Sprintf("auto-sync daily logs - %s", time.Now().Format("2006-01-02 15:04"))
	commitCmd := exec.Command("git", "commit", "-m", commitMsg)
	commitCmd.Dir = dataDir
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to commit logs: %w", err)
	}

	// push to remote upstream tracking branch
	pushCmd := exec.Command("git", "push", "origin", "HEAD")
	pushCmd.Dir = dataDir
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push to GitHub: %w", err)
	}

	fmt.Println("Data successfully pushed.")
	return nil
}
