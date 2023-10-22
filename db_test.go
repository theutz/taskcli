package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetTasksByStatus(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tc.want.Name, tc.want.Project); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			tasks, err := tDB.getTasksByStatus(tc.want.Status)
			if err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			if len(tasks) < 1 {
				t.Fatalf("expected 1 value, got %#v", tasks)
			}
			tc.want.Created = tasks[0].Created
			if !reflect.DeepEqual(tasks[0], tc.want) {
				t.Fatalf("got: %#v, want: %#v", tasks, tc.want)
			}
		})
	}
}

func setup() *taskDB {
	path := filepath.Join(os.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	t := taskDB{db, path}
	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			log.Fatal(err)
		}
	}
	return &t
}

func teardown(tDB *taskDB) {
	tDB.db.Close()
	os.Remove(tDB.dataDir)
}
