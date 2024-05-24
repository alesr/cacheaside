package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ repository = &repoMock{}

type repoMock struct {
	storeFunc func(item Item)
	getFunc   func(id string) (*Item, error)
	allFunc   func() []*Item
}

func (r *repoMock) Store(item Item) {
	r.storeFunc(item)
}

func (r *repoMock) Get(id string) (*Item, error) {
	return r.getFunc(id)
}

func (r *repoMock) All() []*Item {
	return r.allFunc()
}

func TestNew(t *testing.T) {
	t.Parallel()

	// Arrange

	repo := repoMock{}

	// Act

	observed := New(&repo)

	// Assert

	assert.NotNil(t, observed)
	assert.IsType(t, &Service{}, observed)
	assert.Equal(t, &repo, observed.repo)
}

func TestService_Fetch(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		repoResult    func() (*Item, error)
		expectedItem  *Item
		expectedError error
	}{
		{
			name: "should return the item from the repository",
			repoResult: func() (*Item, error) {
				return &Item{ID: "1"}, nil
			},
			expectedItem: &Item{ID: "1"},
		},
		{
			name: "should return an error when the repository fails",
			repoResult: func() (*Item, error) {
				return nil, assert.AnError
			},
			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange

			var repoWasCalled bool
			repo := repoMock{
				getFunc: func(_ string) (*Item, error) {
					repoWasCalled = true
					return tc.repoResult()
				},
			}

			s := New(&repo)

			// Act

			item, err := s.Fetch("1")

			// Assert

			assert.Equal(t, tc.expectedItem, item)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.True(t, repoWasCalled)
		})
	}
}
