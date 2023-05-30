package usermemstorage

import (
	"context"
	"sync"
)

type userMemRepository struct {
	sync.RWMutex
	userStorageCache map[string]map[string]string
}

// AddUserURL add short and original urls pair for user.
func (r *userMemRepository) AddUserURL(c context.Context, userID string, shortURL string, originalURL string) error {
	r.Lock()
	defer r.Unlock()
	r.userStorageCache[userID] = make(map[string]string)
	r.userStorageCache[userID][shortURL] = originalURL
	return nil
}

// GetUserURLsByUserID return all url pairs for user.
func (r *userMemRepository) GetUserURLsByUserID(c context.Context, userID string) (map[string]string, bool) {
	r.RLock()
	defer r.RUnlock()
	urls, ok := r.userStorageCache[userID]
	return urls, ok
}

// New create userStorageCache instance.
func New() *userMemRepository {
	return &userMemRepository{
		userStorageCache: make(map[string]map[string]string),
	}
}
