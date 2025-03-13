package main

import (
	"testing"

	"github.com/turut4/social/internal/store"
	"github.com/turut4/social/internal/store/cache"
	"go.uber.org/zap"
)

func NewTestApplication(t *testing.T) *application {
	t.Helper()

	return &application{
		logger:       zap.Must(zap.NewProduction()).Sugar(),
		store:        store.NewMockStore(t),
		cacheStorage: cache.NewMockCache(t),
	}
}
