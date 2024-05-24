package memdb

import (
	"testing"

	"github.com/alesr/cacheaside/internal/repository"
	"github.com/alesr/cacheaside/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Act

	db := New()

	// Assert

	assert.NotNil(t, db)
	assert.Empty(t, db)
	assert.IsType(t, &MemDB{}, db)
}

func TestMemDB_Get(t *testing.T) {
	t.Parallel()

	// Arrange

	db := New()

	id := "foo-id"
	item := service.Item{ID: id}

	type bar struct {
		ID string
	}

	testCases := []struct {
		name          string
		givenID       string
		storeFunc     func() any
		expectedItem  *service.Item
		expectedError error
	}{
		{
			name:         "found",
			givenID:      "foo-id",
			storeFunc:    func() any { return item },
			expectedItem: &item,
		},
		{
			name:          "not found",
			givenID:       "bar",
			expectedError: repository.ErrItemNotFound,
		},
		{
			name:          "invalid item",
			givenID:       "bar-id",
			storeFunc:     func() any { return bar{} },
			expectedError: repository.ErrInvalidItem,
		},
	}

	// Act

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.storeFunc != nil {
				db.Map.Store(tc.givenID, tc.storeFunc())
			}

			observed, err := db.Get(tc.givenID)

			if tc.expectedError == nil {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedItem, observed)
			} else {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, observed)
			}
		})
	}
}

func TestMemDB_Store(t *testing.T) {
	t.Parallel()

	// Arrange

	db := New()

	id := "foo-id"
	item := service.Item{ID: id}

	// Act

	db.Store(item)

	// Assert

	observed, found := db.Load(id)
	require.True(t, found)

	assert.Equal(t, item, observed)
}

func TestMemDB_All(t *testing.T) {
	t.Parallel()

	// Arrange

	db := New()

	id1 := "foo-id"
	item1 := service.Item{ID: id1}

	id2 := "bar-id"
	item2 := service.Item{ID: id2}

	db.Map.Store(id1, item1)
	db.Map.Store(id2, item2)

	// Act

	observed := db.All()

	// Assert

	assert.Len(t, observed, 2)
	assert.Contains(t, observed, &item1)
	assert.Contains(t, observed, &item2)
}
