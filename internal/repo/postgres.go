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

func (d *DB) AddURL(c context.Context, id entity.URLID, url entity.OriginalURL) {
	insertURL, err := d.db.PrepareContext(c, "INSERT INTO urls (original_url, short_url) VALUES ($1, $2);")
	if err != nil {
		log.Println("ERROR :", err)
	}
	defer insertURL.Close()
	_, err = insertURL.ExecContext(c, string(url), string(id))
	if err != nil {
		log.Println("ERROR :", err)
	}
}

func (d *DB) BatchDeleteShortURLs(c context.Context, urls []entity.ModelURLForDelete) error {
	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()
	updateDeletedURL, err := d.db.PrepareContext(
		ctx, "UPDATE urls SET is_delete = true WHERE user_id = $1 AND short_url = $2;",
	)
	if err != nil {
		return err
	}
	defer updateDeletedURL.Close()
	for _, url := range urls {
		_, err = updateDeletedURL.ExecContext(ctx, url.UserID, url.ShortURL)
		if err != nil {
			log.Println("ERROR :", err)
			return err
		}
	}
	return nil
}

func (d *DB) GetURLByID(c context.Context, id entity.URLID) (entity.OriginalURL, error) {
	var originalURL string
	var isDeleted bool
	selectOriginalURL, err := d.db.PrepareContext(
		c, "SELECT original_url, is_delete  FROM urls WHERE short_url=$1;",
	)
	if err != nil {
		return "", err
	}
	defer selectOriginalURL.Close()

	err = selectOriginalURL.QueryRowContext(c, string(id)).Scan(&originalURL, &isDeleted)
	if isDeleted {
		return "", ErrURLDeleted
	}
	if err != nil {
		log.Println("ERR: ", err)
		return "", err
	}
	return entity.OriginalURL(originalURL), nil
}

func (d *DB) AddUserURL(c context.Context, userID string, shortURL string, originalURL string) error {
	insertUserURL, err := d.db.PrepareContext(
		c, "INSERT INTO urls (user_id, original_url, short_url) VALUES ($1, $2, $3);",
	)
	if err != nil {
		log.Println("ERROR preparing query:", err)
		return err
	}
	defer insertUserURL.Close()
	_, err = insertUserURL.ExecContext(c, userID, originalURL, shortURL)
	if err != nil {
		log.Println("ERROR executing query:", err)
		return ErrURLAlreadyExists
	}
	return nil
}

func (d *DB) GetUserURLsByUserID(c context.Context, userID string) (map[string]string, bool) {
	userURLs := make(map[string]string)
	var modelURL entity.UserURL

	selectUserURLs, err := d.db.PrepareContext(c, "SELECT short_url, original_url FROM urls WHERE user_id=$1;")
	if err != nil {
		log.Println("ERROR :", err)
		return userURLs, false
	}
	defer selectUserURLs.Close()

	row, err := selectUserURLs.QueryContext(c, userID)
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
