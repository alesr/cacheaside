package memcache

import (
	"io"
	"log/slog"
	"testing"

	"github.com/alesr/cacheaside/internal/service"
	"github.com/stretchr/testify/assert"
)

var _ repository = &repoMock{}

type repoMock struct {
	storeFunc func(item service.Item)
	getFunc   func(id string) (*service.Item, error)
	allFunc   func() []*service.Item
}

func (r *repoMock) Store(item service.Item) {
	r.storeFunc(item)
}

func (r *repoMock) Get(id string) (*service.Item, error) {
	return r.getFunc(id)
}

func (r *repoMock) All() []*service.Item {
	return r.allFunc()
}

func TestNew(t *testing.T) {
	t.Parallel()

	// Arrange

	logger := slog.Logger{}
	repo := repoMock{}

	// Act

	observed := New(&logger, &repo)

	// Assert

	assert.NotNil(t, observed)
	assert.IsType(t, &MemCache{}, observed)
	assert.Equal(t, &logger, observed.logger)
	assert.Equal(t, &repo, observed.repository)
	assert.Implements(t, (*repository)(nil), observed)
}

func TestMemCache_Get(t *testing.T) {
	t.Parallel()

	t.Run("should return the item from the cache", func(t *testing.T) {
		t.Parallel()

		// Arrange

		var repoWasCalled bool
		repo := repoMock{
			getFunc: func(_ string) (*service.Item, error) {
				repoWasCalled = true
				return nil, nil
			},
		}

		cache := New(noopLogger(), &repo)

		id := "foo-id"
		item := service.Item{ID: id}

		// Insert the item into the cache.
		cache.cache.Store(id, item)

		// Act

		observed, err := cache.Get(id)

		// Assert

		assert.NoError(t, err)
		assert.Equal(t, &item, observed)
		assert.False(t, repoWasCalled)
	})

	t.Run("should return the item from the repository", func(t *testing.T) {
		t.Parallel()

		// Arrange

		givenID := "foo-id"
		givenItem := service.Item{ID: givenID}

		var repoWasCalled bool
		repo := repoMock{
			getFunc: func(id string) (*service.Item, error) {
				repoWasCalled = true
				assert.Equal(t, id, id)
				return &givenItem, nil
			},
		}

		cache := New(noopLogger(), &repo)

		// Act

		observed, err := cache.Get(givenID)

		// Assert

		assert.NoError(t, err)
		assert.Equal(t, &givenItem, observed)
		assert.True(t, repoWasCalled)
	})

	t.Run("should return the item from the repository if the item in the cache is invalid", func(t *testing.T) {
		t.Parallel()

		// Arrange

		id := "foo-id"
		expected := service.Item{ID: id}

		var repoWasCalled bool
		repo := repoMock{
			getFunc: func(_ string) (*service.Item, error) {
				repoWasCalled = true
				return &expected, nil
			},
		}

		cache := New(noopLogger(), &repo)

		// Insert an invalid item into the cache.
		cache.cache.Store(id, "invalid")

		// Act

		observed, err := cache.Get(id)

		// Assert

		assert.NoError(t, err)
		assert.Equal(t, &expected, observed)
		assert.True(t, repoWasCalled)
	})

	t.Run("should return an error if the item is not in the cache and the repository returns an error", func(t *testing.T) {
		t.Parallel()

		// Arrange

		id := "foo-id"
		expectedErr := assert.AnError

		var repoWasCalled bool
		repo := repoMock{
			getFunc: func(_ string) (*service.Item, error) {
				repoWasCalled = true
				return nil, expectedErr
			},
		}

		cache := New(noopLogger(), &repo)

		// Act

		observed, err := cache.Get(id)

		// Assert

		assert.Nil(t, observed)
		assert.True(t, repoWasCalled)
		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestMemCache_Store(t *testing.T) {
	t.Parallel()

	// Arrange

	givenItem := service.Item{ID: "foo-id"}

	var repoWasCalled bool
	repo := repoMock{
		storeFunc: func(item service.Item) {
			assert.Equal(t, givenItem, item)
			repoWasCalled = true
		},
	}

	cache := New(noopLogger(), &repo)

	// Act

	cache.Store(givenItem)

	// Assert

	_, found := cache.cache.Load(givenItem.ID)
	assert.True(t, found)

	assert.True(t, repoWasCalled)
}

func noopLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}
