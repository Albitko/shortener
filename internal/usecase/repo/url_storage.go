package repo

import "github.com/Albitko/shortener/internal/entity"

var storage = make(map[entity.UrlId]entity.OriginalURL)

type Repository interface {
	AddUrl(entity.UrlId, entity.OriginalURL)
	GetUrlById(entity.UrlId) (entity.OriginalURL, bool)
}

type repository struct{}

func (r *repository) AddUrl(id entity.UrlId, url entity.OriginalURL) {
	storage[id] = url
}

func (r *repository) GetUrlById(id entity.UrlId) (entity.OriginalURL, bool) {
	url, ok := storage[id]
	return url, ok
}

func NewRepository() Repository {
	return &repository{}
}
