// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

/*

Package service provides a simple service on which we can experiment
tests

Part of the Test Double repository I also published
https://github.com/alr-lab/test-double-go

*/
package service

import "context"

type (
	// Service describes a service
	Service struct {
		store Store
	}

	// Store defines a contract for a datastore
	Store interface {
		GetCustomerEmail(ctx context.Context, id int) (string, error)
	}
)

// New returns a new service
func New(store Store) Service {
	return Service{store: store}
}

// Get returns a specific customer email
func (s Service) Get(ctx context.Context) string {
	email, err := s.store.GetCustomerEmail(ctx, 42)
	if err != nil {
		return ""
	}

	return email
}
