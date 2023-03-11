package repo

import (
	"bufio"
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
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0777)
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

func (r *memRepository) AddURL(id entity.URLID, url entity.OriginalURL) {
	r.Lock()
	defer r.Unlock()
	r.storageCache[id] = url

	if r.isFileStorageSet {
		if _, err := r.writer.WriteString(string(id) + "|" + string(url) + "\n"); err != nil {
			return
		}
		err := r.writer.Flush()
		if err != nil {
			return
		}
	}
}

func (r *memRepository) GetURLByID(id entity.URLID) (entity.OriginalURL, bool) {
	r.RLock()
	defer r.RUnlock()
	url, ok := r.storageCache[id]
	return url, ok
}
func (r *memRepository) Close() error {
	return r.fileStorage.Close()
}
