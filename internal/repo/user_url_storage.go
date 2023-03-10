package repo

import (
	"sync"
)

type userMemRepository struct {
	sync.RWMutex
	userStorageCache map[string]map[string]string
}

func (r *userMemRepository) AddUserURL(userID string, shortURL string, originalURL string) {
	r.Lock()
	defer r.Unlock()
	r.userStorageCache[userID] = make(map[string]string)
	r.userStorageCache[userID][shortURL] = originalURL
}

func (r *userMemRepository) GetUserURLsByUserID(userID string) (map[string]string, bool) {
	r.RLock()
	defer r.RUnlock()
	urls, ok := r.userStorageCache[userID]
	return urls, ok
}

func NewUserRepo() *userMemRepository {
	return &userMemRepository{
		userStorageCache: make(map[string]map[string]string),
	}
}
