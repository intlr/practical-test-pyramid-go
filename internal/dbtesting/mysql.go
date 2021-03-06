package dbtesting

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" // needed to use MySQL
	"gopkg.in/testfixtures.v2"
)

const (
	hostname = "datastore"
	username = "root"
	password = "root"
)

// DatabaseHelper provides a connection to a datastore with the application
// schema and some fixtures is passed as an argument
func DatabaseHelper(t *testing.T, fixtureDir string) *sql.DB {
	t.Helper()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		t.Fatalf("Unable to generate uuid")
	}

	schema := fmt.Sprintf("a_%x_%x_%x_%x_%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

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

	if fixtureDir != "" {
		testfixtures.SkipDatabaseNameCheck(true)
		fixtures, err := testfixtures.NewFolder(conn, &testfixtures.MySQL{}, fixtureDir)
		if err != nil {
			t.Fatalf("Unable to create fixtures from fixture directory, err = %s", err)
		}

		if err := fixtures.Load(); err != nil {
			t.Fatalf("Unable to load fixtures, err = %s", err)
		}
	}

	return conn
}
