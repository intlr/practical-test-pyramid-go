// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

package store

import (
	"context"
	"database/sql"
	"fmt"
)

// Store describes a datastore object
type Store struct {
	conn *sql.DB
}

// SetConn attaches the MySQL connection to the store
func (st *Store) SetConn(conn *sql.DB) *Store {
	st.conn = conn
	return st
}

// GetCustomerEmail returns a customer email from the datastore given a
// customer identifier
func (st *Store) GetCustomerEmail(ctx context.Context, id int) (string, error) {
	query := `SELECT email FROM Customers WHERE id = ?`

	args := []interface{}{id}
	res, err := st.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("unable to query database, err = %s", err)
	}
	defer func() { _ = res.Close() }()

	email := ""
	if res.Next() {
		if err := res.Scan(&email); err != nil {
			return email, fmt.Errorf("unable to parse customer email, err = %s", err)
		}
	}

	return email, nil
}
