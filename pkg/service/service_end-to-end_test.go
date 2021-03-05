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
		want     string
		fixtures string
		ctx      context.Context
	}{
		"Successful": {
			want:     "foo",
			fixtures: "successful",
			ctx:      context.Background(),
		},
		"Unexisting": {
			want:     "",
			fixtures: "unexisting",
			ctx:      context.Background(),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			conn := dbtesting.DatabaseHelper(t, fmt.Sprintf("fixtures/%s", tc.fixtures))
			defer conn.Close()
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
