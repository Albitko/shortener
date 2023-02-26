package repo

import (
	"bufio"
	"fmt"
	"github.com/Albitko/shortener/internal/entity"
	"log"
	"os"
	"strings"
)

type fileRepository struct {
	file    *os.File
	writer  *bufio.Writer
	scanner *bufio.Scanner
}

func NewFileRepository(filePath string) *fileRepository {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil
	}
	return &fileRepository{
		file:    file,
		writer:  bufio.NewWriter(file),
		scanner: bufio.NewScanner(file),
	}
}

func (r *fileRepository) AddURL(id entity.URLID, url entity.OriginalURL) {

	if _, err := r.writer.WriteString(string(id) + "|" + string(url) + "\n"); err != nil {
		return
	}
	fmt.Println("WRITE " + string(id) + "|" + string(url))
	err := r.writer.Flush()
	if err != nil {
		return
	}
}

func (r *fileRepository) GetURLByID(id entity.URLID) (entity.OriginalURL, bool) {

	for r.scanner.Scan() {
		record := strings.Split(r.scanner.Text(), "|")
		if record[0] == string(id) {
			return entity.OriginalURL(record[1]), true
		}
	}

	if err := r.scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return "", false
}

func (r *fileRepository) Close() error {
	return r.file.Close()
}
