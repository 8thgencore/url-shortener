package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of the Storage type, initializing the SQLite database.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Close closes the SQLite database connection.
func (s *Storage) Close() error {
	const op = "storage.sqlite.Close"

	// Close the database connection
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// AliasExists checks whether the specified alias exists in the database.
func (s *Storage) AliasExists(alias string) (bool, error) {
	const op = "storage.sqlite.AliasExists"

	stmt, err := s.db.Prepare(`SELECT COUNT(*) FROM url WHERE alias = ?`)
	if err != nil {
		return false, fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var count int
	err = stmt.QueryRow(alias).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return count > 0, nil
}

// URLExists checks whether the specified URL exists in the database.
func (s *Storage) URLExists(urlToCheck string) (bool, error) {
	const op = "storage.sqlite.URLExists"

	stmt, err := s.db.Prepare(`SELECT COUNT(*) FROM url WHERE url = ?`)
	if err != nil {
		return false, fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var count int
	err = stmt.QueryRow(urlToCheck).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return count > 0, nil
}

// GetAliasByURL retrieves the alias associated with a given URL from the database.
func (s *Storage) GetAliasByURL(urlToFind string) (string, error) {
	const op = "storage.sqlite.GetAliasByURL"

	stmt, err := s.db.Prepare(`SELECT alias FROM url WHERE url = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var alias string
	err = stmt.QueryRow(urlToFind).Scan(&alias)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return alias, nil
}

// SaveURL adds a new URL and alias to the database.
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		// TODO: refactoring this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

// GetURL retrieves the URL associated with a given alias from the database.
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

// DeleteURL removes a URL and its associated alias from the database.
func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}
