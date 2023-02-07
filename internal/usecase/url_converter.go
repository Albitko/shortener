package usecase

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

type UrlConverter interface {
	UrlToId(url entity.OriginalURL) entity.UrlId
	IdToUrl(entity.UrlId) (entity.OriginalURL, bool)
}

type urlConverter struct {
	repo repo.Repository
}

func (uc *urlConverter) UrlToId(url entity.OriginalURL) entity.UrlId {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	id := entity.UrlId(base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10])
	uc.repo.AddUrl(id, url)
	return id
}

func (uc *urlConverter) IdToUrl(id entity.UrlId) (entity.OriginalURL, bool) {
	url, ok := uc.repo.GetUrlById(id)
	return url, ok
}

func NewUrlConverter(r repo.Repository) UrlConverter {
	return &urlConverter{
		repo: r,
	}
}
