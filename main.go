package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	gap "github.com/muesli/go-app-paths"
)

func setupPath() string {
	scope := gap.NewScope(gap.User, "tasks")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}
	var taskDir string
	if len(dirs) > 0 {
		taskDir = dirs[0]
	} else {
		taskDir, _ = os.UserHomeDir()
	}
	if err := initTaskDir(taskDir); err != nil {
		log.Fatal(err)
	}
	return taskDir
}

func openDB(path string) (*taskDB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(path, "tasks.db"))
	if err != nil {
		return nil, err
	}
	t := taskDB{db, path}
	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
