package sqlstore

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			b := &bytes.Buffer{}
			fmt.Fprintf(b, "truncate %s cascade", strings.Join(tables, ", "))
			if _, err := db.Exec(b.String()); err!=nil{
				log.Fatal(err, b.String())
			}
		}

		db.Close()
	}
}
