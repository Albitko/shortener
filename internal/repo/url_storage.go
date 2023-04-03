package repo

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Albitko/shortener/internal/entity"
)

type memRepository struct {
	sync.RWMutex
	storageCache     map[entity.URLID]entity.OriginalURL
	fileStorage      *os.File
	writer           *bufio.Writer
	isFileStorageSet bool
}

func NewRepository(path string) *memRepository {
	dataFromFile := make(map[entity.URLID]entity.OriginalURL)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0o777)
	isFileSet := false

	if path != "" {
		isFileSet = true

		if err != nil {
			return nil
		}

		if stat, _ := file.Stat(); stat.Size() != 0 {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				record := strings.Split(scanner.Text(), "|")
				dataFromFile[entity.URLID(record[0])] = entity.OriginalURL(record[1])
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return &memRepository{
		storageCache:     dataFromFile,
		fileStorage:      file,
		writer:           bufio.NewWriter(file),
		isFileStorageSet: isFileSet,
	}
}

func (r *memRepository) AddURL(c context.Context, id entity.URLID, url entity.OriginalURL) {
	r.Lock()
	defer r.Unlock()
	r.storageCache[id] = url

	if r.isFileStorageSet {
		if _, err := r.writer.WriteString(string(id) + "|" + string(url) + "\n"); err != nil {
			log.Print("AddURL failed write string: %w", err)
		}
		err := r.writer.Flush()
		if err != nil {
			log.Print("AddURL failed flush: %w", err)
		}
	}
}

func (r *memRepository) BatchDeleteShortURLs(c context.Context, ids []entity.ModelURLForDelete) error {
	r.Lock()
	defer r.Unlock()
	for _, id := range ids {
		r.storageCache[entity.URLID(id.ShortURL)] = ""
	}
	return nil
}

func (r *memRepository) GetURLByID(c context.Context, id entity.URLID) (entity.OriginalURL, error) {
	r.RLock()
	defer r.RUnlock()
	var err error
	url, ok := r.storageCache[id]
	switch {
	case ok && string(url) == "":
		err = ErrURLDeleted
	case ok:
		err = nil
	default:
		err = errors.New("no needed value in map")
	}
	return url, err
}

func (r *memRepository) Close() error {
	return r.fileStorage.Close()
}
