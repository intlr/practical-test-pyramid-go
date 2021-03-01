package store_test

import (
	"context"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/store"
)

func TestStore(t *testing.T) {
	// Arrange
	conn := dbtesting.DatabaseHelper(t, "")
	defer func() { _ = conn.Close() }()
	st := &store.Store{}
	st.SetConn(conn)
	tt := map[string]struct {
		id   int
		want string
	}{
		"Valid customer identifier will return valid email": {
			id:   42,
			want: "fake",
		},
		"Invalid customer identifier will return empty string": {
			id:   1337,
			want: "",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
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
