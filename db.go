package main

import (
	"database/sql"
	"os"
)

type taskDB struct {
	db      *sql.DB
	dataDir string
}

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o700)
		}
		return err
	}
	return nil
}
