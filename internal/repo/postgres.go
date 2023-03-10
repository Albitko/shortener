package repo

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DB struct {
	db  *sql.DB
	ctx context.Context
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) Ping() error {
	ctx, cancel := context.WithTimeout(d.ctx, 1*time.Second)
	defer cancel()
	err := d.db.PingContext(ctx)
	if err != nil {
		log.Print("ERROR: ", err, "\n")
	}
	return err
}

func NewPostgres(ctx context.Context, psqlConn string) *DB {
	db, err := sql.Open("pgx", psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &DB{
		db:  db,
		ctx: ctx,
	}
}
