// +build integration

package store_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/store"
)

func TestStore(t *testing.T) {
	// Arrange
	tt := map[string]struct {
		id      int
		want    string
		fixture string
	}{
		"Valid customer identifier will return valid email": {
			id:      42,
			want:    "fake",
			fixture: "successful",
		},
		"Invalid customer identifier will return empty string": {
			id:      1337,
			want:    "",
			fixture: "unexisting",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			fixtureDir := fmt.Sprintf("testdata/%s", tc.fixture)
			conn := dbtesting.DatabaseHelper(t, fixtureDir)
			defer func() { _ = conn.Close() }()
			st := (&store.Store{}).SetConn(conn)

			// Act
			got, err := st.GetCustomerEmail(context.Background(), tc.id)

			// Assert
			if err != nil {
				t.Fatalf("Unable to get customer email, err = %s", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
