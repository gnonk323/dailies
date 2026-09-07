package cmd

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"dailies/pkg/server"
	"dailies/pkg/storage"

	"github.com/spf13/cobra"
)

var listenAddr string
var openBrowser bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the dailies HTTP server (owns the SQLite database)",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := storage.DefaultDBPath()
		store, err := storage.OpenSQLite(dbPath)
		if err != nil {
			log.Fatalf("Failed to open database at %s: %v", dbPath, err)
		}
		defer store.Close()
		log.Printf("Using database %s", dbPath)

		if openBrowser {
			go openURL("http://" + listenAddr)
		}

		if err := server.Start(listenAddr, store); err != nil {
			log.Fatalf("Server unexpectedly crashed: %v", err)
		}
	},
}

func init() {
	defaultListen := os.Getenv("DAILIES_LISTEN")
	if defaultListen == "" {
		defaultListen = "127.0.0.1:8080"
	}
	serveCmd.Flags().StringVarP(&listenAddr, "listen", "l", defaultListen, "Address to bind (use a tailnet IP for VPS access)")
	serveCmd.Flags().BoolVarP(&openBrowser, "open", "o", false, "Open the dashboard automatically in the default browser")
	RootCmd.AddCommand(serveCmd)
}

func openURL(url string) {
	time.Sleep(100 * time.Millisecond)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}

	_ = cmd.Start()
}
