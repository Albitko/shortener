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

type userRepository interface {
	AddUserURL(userID string, shortURL string, originalURL string)
	GetUserURLsByUserID(userID string) (map[string]string, bool)
}

type urlConverter struct {
	repo     repository
	userRepo userRepository
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

func (uc *urlConverter) UserIDToURLs(userID string) (map[string]string, bool) {
	urls, ok := uc.userRepo.GetUserURLsByUserID(userID)
	return urls, ok
}

func (uc *urlConverter) AddUserURL(userID string, shortURL string, originalURL string) {
	uc.userRepo.AddUserURL(userID, shortURL, originalURL)
}

func NewURLConverter(r repository, u userRepository) *urlConverter {
	return &urlConverter{
		repo:     r,
		userRepo: u,
	}
}
