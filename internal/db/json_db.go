package db

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")

type Schema struct {
	Appointments []Appointment `json:"appointments"`
	Users        []User        `json:"users"`
}

type Database struct {
	filePath string
	mu       sync.Mutex
}

func NewDatabase(path string) *Database {
	return &Database{filePath: path}
}

func (db *Database) Save(data Schema) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.save(data)
}

func (db *Database) save(data Schema) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.filePath, bytes, 0644)
}

func (db *Database) read() (Schema, error) {
	var data Schema

	file, err := os.ReadFile(db.filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Schema{
				Appointments: []Appointment{},
				Users:        []User{},
			}, nil
		}
		return data, err
	}

	if len(file) == 0 {
		return data, nil
	}

	err = json.Unmarshal(file, &data)
	return data, err
}

func (db *Database) Read() (Schema, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.read()
}
