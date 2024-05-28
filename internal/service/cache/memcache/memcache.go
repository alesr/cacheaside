package memcache

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/alesr/cacheaside/internal/service"
)

type repository interface {
	Store(item service.Item)
	Get(id string) (*service.Item, error)
	All() []*service.Item
}

// MemCache is a cache implementation that stores items in memory.
// It is uses as a cache-aside pattern.
type MemCache struct {
	logger *slog.Logger
	cache  sync.Map
	repository
}

// New creates a new dbCache instance.
func New(logger *slog.Logger, repo repository) *MemCache {
	return &MemCache{
		logger:     logger,
		cache:      sync.Map{},
		repository: repo,
	}
}

// Get retrieves an item by its ID from the cache or the database.
func (c *MemCache) Get(id string) (*service.Item, error) {
	// It first checks the cache.

	v, found := c.cache.Load(id)
	if found {
		c.logger.Info("cache hit", slog.String("id", id))

		item, ok := v.(service.Item)
		if ok {
			return &item, nil
		}

		c.logger.Error("Invalid item in cache. Falling back to database and invalidating cache.", slog.String("id", id))
		c.cache.Delete(id)
	} else {
		c.logger.Info("cache miss", slog.String("id", id))
	}

	// If the item is not in the cache, it retrieves it from the database.
	item, err := c.repository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("could not get item from database: %w", err)
	}

	// And populates the cache.
	c.cache.Store(item.ID, *item)

	return item, nil
}

// Store inserts a new item into the database and the cache.
func (c *MemCache) Store(item service.Item) {
	c.repository.Store(item)
	c.cache.Store(item.ID, item)
}
