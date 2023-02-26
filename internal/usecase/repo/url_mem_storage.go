package repo

import (
	"sync"

	"github.com/Albitko/shortener/internal/entity"
)

type memRepository struct {
	sync.RWMutex
	storage map[entity.URLID]entity.OriginalURL
}

func NewMemRepository() *memRepository {
	return &memRepository{
		storage: make(map[entity.URLID]entity.OriginalURL),
	}
}

func (r *memRepository) AddURL(id entity.URLID, url entity.OriginalURL) {
	r.Lock()
	defer r.Unlock()
	r.storage[id] = url
}

func (r *memRepository) GetURLByID(id entity.URLID) (entity.OriginalURL, bool) {
	r.RLock()
	defer r.RUnlock()
	url, ok := r.storage[id]
	return url, ok
}
