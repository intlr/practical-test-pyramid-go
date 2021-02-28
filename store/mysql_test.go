package store_test

import (
	"context"
	"testing"

	"github.com/alr-lab/ptp/internal/dbtesting"
	"github.com/alr-lab/ptp/store"
)

func TestStore(t *testing.T) {
	want := "fake"

	conn := dbtesting.DatabaseHelper(t)
	defer func(){ _ = conn.Close() }()

	st := &store.Store{}
	st.SetConn(conn)

	got, err := st.GetCustomerEmail(context.Background(), 42)
	if err != nil {
		t.Fatalf("Unable to get customer email, err = %s", err)
	}

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
