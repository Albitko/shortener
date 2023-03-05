package usecase

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/Albitko/shortener/internal/entity"
)

type repository interface {
	AddURL(entity.URLID, entity.OriginalURL)
	GetURLByID(entity.URLID) (entity.OriginalURL, bool)
}

type urlConverter struct {
	repo repository
}

func (uc *urlConverter) URLToID(url entity.OriginalURL) entity.URLID {
	hasher := sha1.New()
	hasher.Write([]byte(url))
	id := entity.URLID(base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10])
	uc.repo.AddURL(id, url)
	return id
}

func (uc *urlConverter) IDToURL(id entity.URLID) (entity.OriginalURL, bool) {
	url, ok := uc.repo.GetURLByID(id)
	return url, ok
}

func NewURLConverter(r repository) *urlConverter {
	return &urlConverter{
		repo: r,
	}
}
