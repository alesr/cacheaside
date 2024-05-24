package memdb

import (
	"sync"

	"github.com/alesr/cacheaside/internal/repository"
	"github.com/alesr/cacheaside/internal/service"
)

// MemDB is an in-memory database implementation.
type MemDB struct {
	sync.Map
}

// New creates a new in-memory database instance.
func New() *MemDB {
	return &MemDB{}
}

// Store inserts a new item into the database.
func (db *MemDB) Store(item service.Item) {
	db.Map.Store(item.ID, item)
}

// Get retrieves an item from the database.
func (db *MemDB) Get(id string) (*service.Item, error) {
	v, found := db.Load(id)
	if !found {
		return nil, repository.ErrItemNotFound
	}

	item, ok := v.(service.Item)
	if !ok {
		return nil, repository.ErrInvalidItem
	}
	return &item, nil
}

// All retrieves all items from the database.
func (db *MemDB) All() []*service.Item {
	items := make([]*service.Item, 0)

	db.Range(func(_, value any) bool {
		item, ok := value.(service.Item)
		if !ok {
			return false
		}

		items = append(items, &item)
		return true
	})
	return items
}
