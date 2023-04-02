package workers

import (
	"fmt"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
)

type Task struct {
	UserID       string
	IDsForDelete []string
}

type Queue struct {
	ch chan *Task
}

func NewQueue() *Queue {
	return &Queue{
		ch: make(chan *Task, 1),
	}
}

func (q *Queue) Push(t *Task) {
	q.ch <- t
}

func (q *Queue) PopWait() *Task {
	return <-q.ch
}

type Deleter struct {
	repository *repo.DB
}

func NewDeleter(r *repo.DB) *Deleter {
	return &Deleter{repository: r}
}

func (r *Deleter) Delete(URLsForDelete []entity.ModelURLForDelete) error {
	return r.repository.BatchDeleteShortURLs(URLsForDelete)
}

type Worker struct {
	id      int
	queue   *Queue
	deleter *Deleter
}

func NewWorker(id int, queue *Queue, deleter *Deleter) *Worker {
	w := Worker{
		id:      id,
		queue:   queue,
		deleter: deleter,
	}
	return &w
}

func (w *Worker) Loop() {
	for {
		t := w.queue.PopWait()
		var URLsForDelete []entity.ModelURLForDelete
		var URLForDelete entity.ModelURLForDelete
		for _, url := range t.IDsForDelete {
			URLForDelete.UserID = t.UserID
			URLForDelete.ShortURL = url
			URLsForDelete = append(URLsForDelete, URLForDelete)
		}
		err := w.deleter.Delete(URLsForDelete)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
	}
}
