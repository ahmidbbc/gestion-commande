// Package store is a repository: data access only, ZERO business logic.
package store

// Store is the concrete repository implementation.
type Store struct{}

// New builds an empty Store.
func New() *Store {
	return &Store{}
}

// Ping reports repository readiness.
func (s *Store) Ping() error {
	return nil
}
