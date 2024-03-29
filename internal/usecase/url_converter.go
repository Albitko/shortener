package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo/postgres"
)

type repository interface {
	AddURL(context.Context, entity.URLID, entity.OriginalURL)
	GetURLByID(context.Context, entity.URLID) (entity.OriginalURL, error)
	BatchDeleteShortURLs(context.Context, []entity.ModelURLForDelete) error
	GetURLsCount(c context.Context) (int, error)
}

type userRepository interface {
	AddUserURL(c context.Context, userID string, shortURL string, originalURL string) error
	GetUserURLsByUserID(c context.Context, userID string) (map[string]string, bool)
	GetUsersCount(c context.Context) (int, error)
}

type urlConverter struct {
	repo     repository
	userRepo userRepository
	pg       *postgres.DB
}

// GetStats return urls and users in service
func (uc *urlConverter) GetStats(ctx context.Context) (entity.URLStats, error) {
	var stats entity.URLStats

	usersCount, err := uc.userRepo.GetUsersCount(ctx)
	if err != nil {
		return stats, err
	}
	urlsCount, err := uc.repo.GetURLsCount(ctx)
	if err != nil {
		return stats, err
	}
	stats.URLsCount = urlsCount
	stats.UsersCount = usersCount

	return stats, nil
}

// URLToID generate short URL and add it to DB.
func (uc *urlConverter) URLToID(ctx context.Context, url entity.OriginalURL, userID string) (entity.URLID, error) {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	id := entity.URLID(base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:10])
	err := uc.userRepo.AddUserURL(ctx, userID, string(id), string(url))
	uc.repo.AddURL(ctx, id, url)
	return id, err
}

// IDToURL get original URL for shorten. Pass task to the next layer.
func (uc *urlConverter) IDToURL(ctx context.Context, id entity.URLID) (entity.OriginalURL, error) {
	url, err := uc.repo.GetURLByID(ctx, id)
	return url, err
}

// BatchDeleteURL prepare data for batch delete in DB and pass data to it.
func (uc *urlConverter) BatchDeleteURL(c context.Context, userID string, shortURLs []string) {
	URLsForDelete := make([]entity.ModelURLForDelete, len(shortURLs))

	var URLForDelete entity.ModelURLForDelete
	for i, url := range shortURLs {
		URLForDelete.UserID = userID
		URLForDelete.ShortURL = url
		URLsForDelete[i] = URLForDelete
	}
	err := uc.repo.BatchDeleteShortURLs(c, URLsForDelete)
	if err != nil {
		log.Println("ERROR update delete URLs ", err)
	}
}

// UserIDToURLs return all user urls.
func (uc *urlConverter) UserIDToURLs(ctx context.Context, userID string) (map[string]string, bool) {
	urls, ok := uc.userRepo.GetUserURLsByUserID(ctx, userID)
	return urls, ok
}

// AddUserURL add shorten url for user.
func (uc *urlConverter) AddUserURL(c context.Context, userID string, shortURL string, originalURL string) error {
	err := uc.userRepo.AddUserURL(c, userID, shortURL, originalURL)
	return err
}

// PingDB check DB connection.
func (uc *urlConverter) PingDB() error {
	err := uc.pg.Ping()
	return err
}

// New create urlConverter instance.
func New(r repository, u userRepository, d *postgres.DB) *urlConverter {
	return &urlConverter{
		repo:     r,
		userRepo: u,
		pg:       d,
	}
}
