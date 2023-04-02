package workers

import (
	"fmt"
	"runtime"

	"github.com/Albitko/shortener/internal/entity"
)

type repository interface {
	BatchDeleteShortURLs([]entity.ModelURLForDelete) error
}
type Task struct {
	UserID       string
	IDsForDelete []string
}

type Queue struct {
	ch chan *Task
}

func newQueue() *Queue {
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
	repo repository
}

func newDeleter(r repository) *Deleter {
	return &Deleter{repo: r}
}

func (r *Deleter) Delete(urlsForDelete []entity.ModelURLForDelete) error {
	return r.repo.BatchDeleteShortURLs(urlsForDelete)
}

type Worker struct {
	id      int
	queue   *Queue
	deleter *Deleter
}

func newWorker(id int, queue *Queue, deleter *Deleter) *Worker {
	w := Worker{
		id:      id,
		queue:   queue,
		deleter: deleter,
	}
	return &w
}

func (w *Worker) loop() {
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

func InitWorkers(r repository) *Queue {
	queue := newQueue()
	wrkrs := make([]*Worker, 0, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wrkrs = append(wrkrs, newWorker(i, queue, newDeleter(r)))
	}

	for _, w := range wrkrs {
		go w.loop()
	}
	return queue
}
