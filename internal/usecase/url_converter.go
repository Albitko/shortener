package usecase

import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
)

type repository interface {
	AddURL(entity.URLID, entity.OriginalURL)
	GetURLByID(entity.URLID) (entity.OriginalURL, error)
	BatchDeleteShortURLs([]entity.ModelURLForDelete) error
}

type userRepository interface {
	AddUserURL(userID string, shortURL string, originalURL string) error
	GetUserURLsByUserID(userID string) (map[string]string, bool)
}

type urlConverter struct {
	repo     repository
	userRepo userRepository
	pg       *repo.DB
}

func (uc *urlConverter) URLToID(url entity.OriginalURL, userID string) (entity.URLID, error) {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	id := entity.URLID(base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10])
	err := uc.userRepo.AddUserURL(userID, string(id), string(url))
	uc.repo.AddURL(id, url)
	return id, err
}

func (uc *urlConverter) IDToURL(id entity.URLID) (entity.OriginalURL, error) {
	url, err := uc.repo.GetURLByID(id)
	return url, err
}

func (uc *urlConverter) BatchDeleteURL(userID string, shortURLs []string) {
	URLsForDelete := make([]entity.ModelURLForDelete, len(shortURLs))

	var URLForDelete entity.ModelURLForDelete
	for i, url := range shortURLs {
		URLForDelete.UserID = userID
		URLForDelete.ShortURL = url
		URLsForDelete[i] = URLForDelete
	}
	err := uc.repo.BatchDeleteShortURLs(URLsForDelete)
	if err != nil {
		log.Println("ERROR update delete URLs ", err)
	}
}

func (uc *urlConverter) UserIDToURLs(userID string) (map[string]string, bool) {
	urls, ok := uc.userRepo.GetUserURLsByUserID(userID)
	return urls, ok
}

func (uc *urlConverter) AddUserURL(userID string, shortURL string, originalURL string) error {
	err := uc.userRepo.AddUserURL(userID, shortURL, originalURL)
	return err
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
