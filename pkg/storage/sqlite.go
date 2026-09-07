package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"dailies/pkg/types"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func OpenSQLite(path string) (*SQLiteStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA busy_timeout=5000;`); err != nil {
		db.Close()
		return nil, err
	}

	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			date TEXT PRIMARY KEY,
			data TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS config (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			data TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	`)
	return err
}

func (s *SQLiteStore) LoadConfig() (types.DailiesConfig, error) {
	var raw string
	err := s.db.QueryRow(`SELECT data FROM config WHERE id = 1`).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return DefaultConfig(), err
	}

	var cfg types.DailiesConfig
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return DefaultConfig(), err
	}
	if cfg.Integrations == nil {
		cfg.Integrations = make(map[string]map[string]bool)
	}
	return cfg, nil
}

func (s *SQLiteStore) SaveConfig(cfg types.DailiesConfig) error {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, err = s.db.Exec(`
		INSERT INTO config (id, data, updated_at)
		VALUES (1, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			data = excluded.data,
			updated_at = excluded.updated_at
	`, string(bytes), now)
	return err
}

func (s *SQLiteStore) LoadEntry(dateStr string) (*types.DailyEntry, error) {
	var raw string
	err := s.db.QueryRow(`SELECT data FROM entries WHERE date = ?`, dateStr).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var entry types.DailyEntry
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *SQLiteStore) SaveEntry(entry *types.DailyEntry) error {
	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, err = s.db.Exec(`
		INSERT INTO entries (date, data, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(date) DO UPDATE SET
			data = excluded.data,
			updated_at = excluded.updated_at
	`, entry.Date, string(bytes), now, now)
	return err
}

func (s *SQLiteStore) ListEntries() ([]types.DailyEntry, error) {
	rows, err := s.db.Query(`SELECT data FROM entries ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []types.DailyEntry{}
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		var entry types.DailyEntry
		if err := json.Unmarshal([]byte(raw), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
