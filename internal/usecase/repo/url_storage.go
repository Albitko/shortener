package repo

import (
	"sync"

	"github.com/Albitko/shortener/internal/entity"
)

type repository struct {
	sync.RWMutex
	storage map[entity.URLID]entity.OriginalURL
}

func NewRepository() *repository {
	return &repository{
		storage: make(map[entity.URLID]entity.OriginalURL),
	}
}

func (r *repository) AddURL(id entity.URLID, url entity.OriginalURL) {
	r.Lock()
	defer r.Unlock()
	r.storage[id] = url
}

func (r *repository) GetURLByID(id entity.URLID) (entity.OriginalURL, bool) {
	r.RLock()
	defer r.RUnlock()
	url, ok := r.storage[id]
	return url, ok
}
