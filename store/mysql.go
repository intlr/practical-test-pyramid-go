package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	conn *sql.DB
}

func (st *Store) SetConn(conn *sql.DB) {
	st.conn = conn
}

func (s Store) GetCustomerEmail(ctx context.Context, id int) (string, error) {
	query := `SELECT email FROM Customers WHERE id = ?`

	args := []interface{}{id}
	res, err := s.conn.QueryContext(ctx, query, args...)
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
