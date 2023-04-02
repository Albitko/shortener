package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Albitko/shortener/internal/entity"
)

const schema = `
 	CREATE TABLE IF NOT EXISTS urls (
 		id serial primary key,
 		user_id text,
 		original_url text not null unique,
 		short_url text not null,
 		is_delete boolean DEFAULT FALSE
 	);
 	`

var (
	ErrURLAlreadyExists = errors.New("URL already exists")
	ErrURLDeleted       = errors.New("URL deleted")
)

type DB struct {
	db  *sql.DB
	ctx context.Context
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) AddURL(id entity.URLID, url entity.OriginalURL) {
	insertURL, err := d.db.Prepare("INSERT INTO urls (original_url, short_url) VALUES ($1, $2);")
	if err != nil {
		log.Println("ERROR :", err)
	}
	defer insertURL.Close()
	_, err = insertURL.Exec(string(url), string(id))
	if err != nil {
		log.Println("ERROR :", err)
		//return ErrURLAlreadyExists
	}
}
func (d *DB) BatchDeleteShortURLs(urls []entity.ModelURLForDelete) error {
	updateDeletedURL, err := d.db.Prepare("UPDATE urls SET is_delete = true WHERE user_id = $1 AND short_url = $2;")
	if err != nil {
		return err
	}
	defer updateDeletedURL.Close()
	for _, url := range urls {
		_, err = updateDeletedURL.Exec(url.UserID, url.ShortURL)
		if err != nil {
			log.Println("ERROR :", err)
			return err
		}
	}
	return nil
}

func (d *DB) GetURLByID(id entity.URLID) (entity.OriginalURL, error) {
	var originalURL string
	var isDeleted bool
	selectOriginalURL, err := d.db.Prepare("SELECT original_url, is_delete  FROM urls WHERE short_url=$1;")
	if err != nil {
		return "", err
	}
	defer selectOriginalURL.Close()

	err = selectOriginalURL.QueryRow(string(id)).Scan(&originalURL, &isDeleted)
	if isDeleted {
		return "", ErrURLDeleted
	}
	if err != nil {
		log.Println("ERR: ", err)
		return "", err
	}
	return entity.OriginalURL(originalURL), nil
}

func (d *DB) AddUserURL(userID string, shortURL string, originalURL string) error {
	insertUserURL, err := d.db.Prepare("INSERT INTO urls (user_id, original_url, short_url) VALUES ($1, $2, $3);")
	if err != nil {
		log.Println("ERROR 1:", err)
		return err
	}
	defer insertUserURL.Close()
	_, err = insertUserURL.Exec(userID, originalURL, shortURL)
	if err != nil {
		log.Println("ERROR 2:", err)
		return ErrURLAlreadyExists
	}
	return nil
}

func (d *DB) GetUserURLsByUserID(userID string) (map[string]string, bool) {
	userURLs := make(map[string]string)
	var modelURL entity.UserURL

	selectUserURLs, err := d.db.Prepare("SELECT short_url, original_url FROM urls WHERE user_id=$1;")
	if err != nil {
		log.Println("ERROR :", err)
		return userURLs, false
	}
	defer selectUserURLs.Close()

	row, err := selectUserURLs.Query(userID)
	if err != nil {
		log.Println("ERROR :", err)
		return userURLs, false
	}
	defer row.Close()

	if err = row.Err(); err != nil {
		log.Println(err)
		return userURLs, false
	}

	for row.Next() {
		err := row.Scan(&modelURL.ShortURL, &modelURL.OriginalURL)
		if err != nil {
			log.Println("ERROR :", err)
			return userURLs, false
		}
		userURLs[modelURL.ShortURL] = modelURL.OriginalURL
	}
	return userURLs, true
}

func (d *DB) Ping() error {
	ctx, cancel := context.WithTimeout(d.ctx, 1*time.Second)
	defer cancel()
	err := d.db.PingContext(ctx)
	if err != nil {
		log.Print("ERROR: ", err, "\n")
		return fmt.Errorf("PingContext failed: %w", err)
	}
	return nil
}

func NewPostgres(ctx context.Context, psqlConn string) *DB {
	db, err := sql.Open("pgx", psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(schema); err != nil {
		log.Fatal(err)
	}
	return &DB{
		db:  db,
		ctx: ctx,
	}
}
