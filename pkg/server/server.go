package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"sort"
	"dailies/pkg/integrations"
	"dailies/pkg/storage"
	"dailies/pkg/types"
)

//go:embed dist
var frontendFS embed.FS

// enableCors sets up simple headers for local React development
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Start spins up the REST server on the specified port
func Start(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/config", handleGetConfig)
	mux.HandleFunc("GET /api/entries", handleGetAllEntries)
	mux.HandleFunc("GET /api/entries/{date}", handleGetEntry)
	mux.HandleFunc("POST /api/entries", handleSaveEntry)
	mux.HandleFunc("POST /api/entries/{date}/fetch/{integration}", handleTriggerIntegration)

	strippedFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(strippedFS))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := strippedFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// if the file doesn't exist, serve index.html
		indexFile, err := strippedFS.Open("index.html")
		if err != nil {
			http.Error(w, "Frontend build not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		seeker, ok := indexFile.(io.ReadSeeker)
		if !ok {
			http.Error(w, "Internal server error: file not seekable", http.StatusInternalServerError)
			return
		}

		stat, _ := indexFile.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), seeker)
	})

	globalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		mux.ServeHTTP(w, r)
	})

	println("Server starting on http://localhost:" + port)
	return http.ListenAndServe(":"+port, globalHandler)
}

// GET /api/config
func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := storage.LoadConfig()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

// GET /api/entries
func handleGetAllEntries(w http.ResponseWriter, r *http.Request) {
	dir := storage.GetDataDirectory()
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Unable to read data directory", http.StatusInternalServerError)
		return
	}

	var allEntries []types.DailyEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			dateStr := strings.TrimSuffix(file.Name(), ".json")
			entry, err := storage.LoadEntry(dateStr)
			if err == nil && entry != nil {
				allEntries = append(allEntries, *entry)
			}
		}
	}

	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].Date > allEntries[j].Date
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allEntries)
}

// GET /api/entries/{date}
func handleGetEntry(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date") // works on Go 1.22+ standard routing
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date syntax. Expecting YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	entry, err := storage.LoadEntry(dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if entry == nil {
		http.Error(w, "Log entry not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// POST /api/entries
func handleSaveEntry(w http.ResponseWriter, r *http.Request) {
	var entry types.DailyEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Malformed JSON payload", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", entry.Date); err != nil {
		http.Error(w, "Invalid or missing 'date' key (YYYY-MM-DD format required)", http.StatusBadRequest)
		return
	}

	if err := storage.SaveEntry(&entry); err != nil {
		http.Error(w, "Failed to persist log data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "successfully written"})
}

// POST /api/entries/{date}/fetch/{integration}
func handleTriggerIntegration(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	integrationName := r.PathValue("integration")

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Make sure the requested integration is valid
	if _, ok := integrations.IntegrationRegistry[integrationName]; !ok {
		http.Error(w, "Requested integration payload driver does not exist", http.StatusNotFound)
		return
	}

	// Triggering your module synchronous runner safely
	// Note: RunManualFetch saves internal records automatically via its internal storage triggers
	integrations.RunManualFetch(integrationName, dateStr)

	// Fetch updated document version to yield back cleanly to React front-end
	updatedEntry, err := storage.LoadEntry(dateStr)
	if err != nil {
		http.Error(w, "Integration finished but failed reload pipeline sequence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedEntry)
}
