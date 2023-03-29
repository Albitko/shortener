package usecase

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
)

type repository interface {
	AddURL(entity.URLID, entity.OriginalURL) error
	GetURLByID(entity.URLID) (entity.OriginalURL, bool)
}

type userRepository interface {
	AddUserURL(userID string, shortURL string, originalURL string)
	GetUserURLsByUserID(userID string) (map[string]string, bool)
}

type urlConverter struct {
	repo     repository
	userRepo userRepository
	pg       *repo.DB
}

func (uc *urlConverter) URLToID(url entity.OriginalURL, userID string, baseURL string) (entity.URLID, error) {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	id := entity.URLID(base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10])
	uc.userRepo.AddUserURL(userID, baseURL+string(id), string(url))
	err := uc.repo.AddURL(id, url)
	return id, err
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

func (uc *urlConverter) PingDB() error {
	err := uc.pg.Ping()
	return err
}

func NewURLConverter(r repository, u userRepository, d *repo.DB) *urlConverter {
	return &urlConverter{
		repo:     r,
		userRepo: u,
		pg:       d,
	}
}
