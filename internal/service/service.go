// Package service holds usecases: business orchestration. It depends on
// repository interfaces, never on the api (handler) layer.
package service

// Repository is the data-access port the service depends on. The concrete
// implementation lives in internal/store.
type Repository interface {
	Ping() error
}

// Service orchestrates business logic over its repository dependencies.
type Service struct {
	repo Repository
}

// New builds a Service with its injected repository.
func New(repo Repository) *Service {
	return &Service{repo: repo}
}
