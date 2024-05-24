package service

import (
	"fmt"

	"github.com/google/uuid"
)

// Item represents a sample item domain entity.
type Item struct {
	ID string
}

type repository interface {
	Store(item Item)
	Get(id string) (*Item, error)
	All() []*Item
}

// Service represents a service that manages items.
type Service struct {
	repo repository
}

// New creates a new service instance.
func New(repo repository) *Service {
	return &Service{repo: repo}
}

// Fetch retrieves an item by its ID.
func (s *Service) Fetch(id string) (*Item, error) {
	item, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("could not get item from repository: %w", err)
	}
	return item, nil
}

// List retrieves all items.
func (s *Service) List() []*Item {
	return s.repo.All()
}

// Create creates a new item with a random ID
// and stores it in the repository.
func (s *Service) Create() *Item {
	item := Item{ID: uuid.New().String()}
	s.repo.Store(item)
	return &item
}
