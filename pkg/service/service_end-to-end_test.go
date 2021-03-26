// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

// +build endtoend

package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/store"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/service"
)

func TestService_fixtures(t *testing.T) {
	// Arrange
	tt := map[string]struct {
		want    string
		fixture string
		ctx     context.Context
	}{
		"Successful": {
			want:    "foo",
			fixture: "successful",
			ctx:     context.Background(),
		},
		"Unexisting": {
			want:    "",
			fixture: "unexisting",
			ctx:     context.Background(),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fixtureDir := fmt.Sprintf("testdata/%s", tc.fixture)
			conn, teardown := dbtesting.DatabaseHelper(t, fixtureDir)
			defer teardown()
			st := &store.Store{}
			st.SetConn(conn)
			serv := service.New(st)

			// Act
			got := serv.Get(tc.ctx)

			// Assert
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
