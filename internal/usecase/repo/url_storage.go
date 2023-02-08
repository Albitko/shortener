package repo

import "github.com/Albitko/shortener/internal/entity"

var storage = make(map[entity.URLID]entity.OriginalURL)

type Repository interface {
	AddURL(entity.URLID, entity.OriginalURL)
	GetURLByID(entity.URLID) (entity.OriginalURL, bool)
}

type repository struct{}

func (r *repository) AddURL(id entity.URLID, url entity.OriginalURL) {
	storage[id] = url
}

func (r *repository) GetURLByID(id entity.URLID) (entity.OriginalURL, bool) {
	url, ok := storage[id]
	return url, ok
}

func NewRepository() Repository {
	return &repository{}
}
