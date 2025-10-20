package database

import (
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/BerylCAtieno/string-analyzer/analysis"
	"github.com/BerylCAtieno/string-analyzer/models"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes a SQLite database and creates a table strings if it does not exist
func InitDB() error {

	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite", "./data/strings.db")
	if err != nil {
		return err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS strings (
		id TEXT PRIMARY KEY,
		value TEXT UNIQUE NOT NULL,
		properties TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return DB.Ping()
}

// InsertString inserts a new analyzed string into the database
func InsertString(id, value string, properties map[string]any) error {
	propsJSON, err := json.Marshal(properties)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
		INSERT INTO strings (id, value, properties, created_at)
		VALUES (?, ?, ?, ?)
	`, id, value, string(propsJSON), time.Now().UTC().Format(time.RFC3339))

	return err
}

// GetStringByValue gets a string record by value
func GetStringByValue(value string) (string, string, string, error) {
	row := DB.QueryRow(`SELECT id, value, properties FROM strings WHERE value = ?`, value)

	var id, val, props string
	err := row.Scan(&id, &val, &props)
	return id, val, props, err
}

// DeleteString deletes a string record by value
func DeleteString(value string) error {
	_, err := DB.Exec(`DELETE FROM strings WHERE value = ?`, value)
	return err
}

// GetAllStrings gets all strings
func GetAllStrings() ([]models.StringReponsePayload, error) {
	rows, err := DB.Query(`SELECT id, value, properties, created_at FROM strings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.StringReponsePayload

	for rows.Next() {
		var id, value, propsJSON, createdAt string
		if err := rows.Scan(&id, &value, &propsJSON, &createdAt); err != nil {
			return nil, err
		}

		var props analysis.Properties
		if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
			return nil, err
		}

		results = append(results, models.StringReponsePayload{
			Id:         id,
			Value:      value,
			Properties: props,
			CreatedAt:  createdAt,
		})
	}

	return results, nil
}
