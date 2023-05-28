package workers

import (
	"context"
	"fmt"
	"runtime"

	"github.com/Albitko/shortener/internal/entity"
)

type repository interface {
	BatchDeleteShortURLs(context.Context, []entity.ModelURLForDelete) error
}

// Task type that represents short url that should be deleted for user.
type Task struct {
	UserID       string
	IDsForDelete []string
}

// Queue of tasks.
type Queue struct {
	ch chan *Task
}

func newQueue() *Queue {
	return &Queue{
		ch: make(chan *Task, 1),
	}
}

// Push add task to queue.
func (q *Queue) Push(t *Task) {
	q.ch <- t
}

// PopWait get task from queue.
func (q *Queue) PopWait() *Task {
	return <-q.ch
}

// Deleter type that delete urls.
type Deleter struct {
	repo repository
	ctx  context.Context
}

func newDeleter(ctx context.Context, r repository) *Deleter {
	return &Deleter{repo: r, ctx: ctx}
}

// Delete delete urls from DB.
func (d *Deleter) Delete(urlsForDelete []entity.ModelURLForDelete) error {
	return d.repo.BatchDeleteShortURLs(d.ctx, urlsForDelete)
}

// Worker checks if queue is not empty and delete urls.
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
		for i := range t.IDsForDelete {
			URLForDelete.UserID = t.UserID
			URLForDelete.ShortURL = t.IDsForDelete[i]
			URLsForDelete = append(URLsForDelete, URLForDelete)
		}
		err := w.deleter.Delete(URLsForDelete)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
	}
}

// InitWorkers create queue and workers.
func InitWorkers(ctx context.Context, r repository) *Queue {
	queue := newQueue()
	wrkrs := make([]*Worker, 0, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wrkrs = append(wrkrs, newWorker(i, queue, newDeleter(ctx, r)))
	}

	for i := range wrkrs {
		go wrkrs[i].loop()
	}
	return queue
}
