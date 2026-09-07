package server

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"dailies/pkg/integrations"
	"dailies/pkg/storage"
	"dailies/pkg/types"
)

//go:embed dist
var frontendFS embed.FS

type Server struct {
	store storage.Store
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func Start(listenAddr string, store storage.Store) error {
	s := &Server{store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/config", s.handleGetConfig)
	mux.HandleFunc("POST /api/config", s.handleSaveConfig)
	mux.HandleFunc("GET /api/entries", s.handleGetAllEntries)
	mux.HandleFunc("GET /api/entries/{date}", s.handleGetEntry)
	mux.HandleFunc("POST /api/entries", s.handleSaveEntry)
	mux.HandleFunc("POST /api/entries/{date}/fetch", s.handleTriggerAllIntegrations)
	mux.HandleFunc("POST /api/entries/{date}/fetch/{integration}", s.handleTriggerIntegration)

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

	log.Printf("Server starting on http://%s", listenAddr)
	return http.ListenAndServe(listenAddr, globalHandler)
}

func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := s.store.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (s *Server) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var cfg types.DailiesConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, "Malformed JSON payload", http.StatusBadRequest)
		return
	}
	if err := s.store.SaveConfig(cfg); err != nil {
		http.Error(w, "Failed to persist config", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}

func (s *Server) handleGetAllEntries(w http.ResponseWriter, r *http.Request) {
	allEntries, err := s.store.ListEntries()
	if err != nil {
		http.Error(w, "Unable to read entries", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, allEntries)
}

func (s *Server) handleGetEntry(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date syntax. Expecting YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	entry, err := s.store.LoadEntry(dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if entry == nil {
		http.Error(w, "Log entry not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleSaveEntry(w http.ResponseWriter, r *http.Request) {
	var entry types.DailyEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Malformed JSON payload", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", entry.Date); err != nil {
		http.Error(w, "Invalid or missing 'date' key (YYYY-MM-DD format required)", http.StatusBadRequest)
		return
	}

	if err := s.store.SaveEntry(&entry); err != nil {
		http.Error(w, "Failed to persist log data", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "successfully written"})
}

func (s *Server) handleTriggerAllIntegrations(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	integrations.RunAllManualFetch(s.store, dateStr)
	updatedEntry, err := s.store.LoadEntry(dateStr)
	if err != nil {
		http.Error(w, "Integration finished but failed reload", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updatedEntry)
}

func (s *Server) handleTriggerIntegration(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	integrationName := r.PathValue("integration")

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	if _, ok := integrations.IntegrationRegistry[integrationName]; !ok {
		http.Error(w, "Requested integration payload driver does not exist", http.StatusNotFound)
		return
	}

	integrations.RunManualFetch(s.store, integrationName, dateStr)

	updatedEntry, err := s.store.LoadEntry(dateStr)
	if err != nil {
		http.Error(w, "Integration finished but failed reload pipeline sequence", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updatedEntry)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
