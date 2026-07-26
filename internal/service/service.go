// Package service holds usecases: business orchestration. It depends on
// repository interfaces, never on the api (handler) layer.
package service

// Order is the core domain entity: a customer order.
type Order struct {
	ID       string
	Customer string
	Total    float64
	Status   string
}

// Repository is the data-access port the service depends on. The concrete
// implementation lives in internal/store.
type Repository interface {
	Ping() error
	Orders() ([]Order, error)
}

// Service orchestrates business logic over its repository dependencies.
type Service struct {
	repo Repository
}

// New builds a Service with its injected repository.
func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListOrders returns all customer orders.
func (s *Service) ListOrders() ([]Order, error) {
	return s.repo.Orders()
}
