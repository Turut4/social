package cache

import (
	"context"
	"testing"

	"github.com/turut4/social/internal/store"
)

func NewMockCache(t *testing.T) Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	return &store.User{}, nil
}

func (m *MockUserStore) Set(ctx context.Context, user *store.User) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) {
}
