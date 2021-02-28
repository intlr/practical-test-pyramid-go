// Package service provides a simple service on which we can experiment
// tests
//
// Part of the Test Double repository I also published
// https://github.com/alr-lab/test-double-go
package service

type (
	// Service describes a service
	Service struct {
		store Store
	}

	// Store defines a contract for a datastore
	Store interface {
		GetCustomerEmail(id int) string
	}
)

// New returns a new service
func New(store Store) Service {
	return Service{store: store}
}

// Get returns a specific customer email
func (s Service) Get() string {
	return s.store.GetCustomerEmail(42)
}
