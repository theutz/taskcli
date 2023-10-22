package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

type task struct {
	ID      uint
	Name    string
	Project string
	Status  string
	Created time.Time
}

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

func (t *taskDB) tableExists(name string) bool {
	if _, err := t.db.Query("SELECT * FROM tasks"); err == nil {
		return true
	}
	return false
}

func (t *taskDB) createTable() error {
	_, err := t.db.Exec(`CREATE TABLE "tasks" ( "id" INTEGER, "name" TEXT NOT NULL, "project" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func (t *taskDB) insert(name, project string) error {
	_, err := t.db.Exec(
		"INSERT INTO tasks(name, project, status, created) VALUES( ?, ?, ?, ?)",
		name,
		project,
		todo.String(),
		time.Now(),
	)
	return err
}

func (t *taskDB) getTasksByStatus(status string) ([]task, error) {
	var tasks []task
	rows, err := t.db.Query("SELECT * FROM tasks WHERE status = ?", status)
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}
