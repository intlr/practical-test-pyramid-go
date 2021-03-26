// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

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

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"gopkg.in/testfixtures.v2"
)

const (
	// Testing database hostname
	hostname = "datastore"

	// Testing database username
	username = "root"

	// Testing database password
	password = "root"
)

// DatabaseHelper provides a connection to a datastore with the application
// schema and some fixtures is passed as an argument
func DatabaseHelper(t *testing.T, fixtureDir string) (*sql.DB, func()) {
	t.Helper()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		t.Fatalf("Unable to generate database schema UUID")
	}

	// Database schema, must start with a letter, ensures the schema is unique
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
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(2 * time.Minute)

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
		fixtures(t, conn, fixtureDir)
	}

	teardown := func() {
		_, _ = conn.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, schema))
		conn.Close()
	}

	return conn, teardown
}

// Run fixtures against newly created database described by the connection
// given as an argument
func fixtures(t *testing.T, conn *sql.DB, fixtureDir string) {
	testfixtures.SkipDatabaseNameCheck(true)

	fixtures, err := testfixtures.NewFolder(conn, &testfixtures.MySQL{}, fixtureDir)
	if err != nil {
		t.Fatalf("Unable to create fixtures from fixture directory, err = %s", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatalf("Unable to load fixtures, err = %s", err)
	}
}
