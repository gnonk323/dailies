package cmd

import (
	"dailies/pkg/server"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var port string
var openBrowser bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Spins up the local REST server backend dashboard mapping engines",
	Run: func(cmd *cobra.Command, args []string) {
		if openBrowser {
			go openURL(fmt.Sprintf("http://localhost:%s", port))
		}

		if err := server.Start(port); err != nil {
			log.Fatalf("Server unexpectedly crashed: %v", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Target running port configuration assignment link")
	serveCmd.Flags().BoolVarP(&openBrowser, "open", "o", false, "Open the dashboard automatically in the default browser")
	
	RootCmd.AddCommand(serveCmd)
}

func openURL(url string) {
	time.Sleep(100 * time.Millisecond)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}

	_ = cmd.Start()
}