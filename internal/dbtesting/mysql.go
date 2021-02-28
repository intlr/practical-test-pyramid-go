package dbtesting

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	schema   = "app"
	hostname = "mysql"
	username = "root"
	password = "root"
)

func DatabaseHelper(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/?%s",
		username,
		password,
		hostname,
		"parseTime=true&readTimeout=1s&writeTimeout=1s&multiStatements=true",
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Unable to open local MySQL connection, err = %s", err)
	}

	conn.SetMaxOpenConns(5)
	conn.SetConnMaxLifetime(10 * time.Minute)

	query := fmt.Sprintf(
		`SET GLOBAL max_connections = 10; DROP DATABASE IF EXISTS %[1]s; CREATE DATABASE %[1]s; USE %[1]s;`,
		schema,
	)

	if _, err := conn.Exec(query); err != nil {
		t.Fatalf("Unable to execute database creation query, err = %s", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	p := path.Join(path.Dir(filename), "skeleton/schema.sql")

	skeleton, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("Unable to read database bootstrap queries, err = %s", err)
	}

	if _, err := conn.Exec(string(skeleton)); err != nil {
		t.Fatalf("Unable to bootstrap database, err = %s", err)
	}

	return conn
}
