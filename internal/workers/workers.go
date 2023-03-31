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
		ch: make(chan *Task, 100),
	}
}

func (q *Queue) Push(t *Task) {
	q.ch <- t
}

func (q *Queue) PopWait() *Task {
	return <-q.ch
}

type Resizer struct {
	repository *repo.DB
}

func NewResizer(r *repo.DB) *Resizer {
	return &Resizer{repository: r}
}

func (r *Resizer) Resize(URLsForDelete []entity.ModelURLForDelete) error {
	//log.Println("DELETING ", URLsForDelete)
	return r.repository.BatchDeleteShortURLs(URLsForDelete)
}

type Worker struct {
	id      int
	queue   *Queue
	resizer *Resizer
}

func NewWorker(id int, queue *Queue, resizer *Resizer) *Worker {
	w := Worker{
		id:      id,
		queue:   queue,
		resizer: resizer,
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
		err := w.resizer.Resize(URLsForDelete)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
	}
}
