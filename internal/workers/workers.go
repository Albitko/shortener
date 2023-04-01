package workers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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

type InputWorker struct {
	ch     chan Task
	done   chan struct{}
	index  int
	ticker *time.Ticker
	ctx    context.Context
	mu     *sync.Mutex
}

type OutputWorker struct {
	id   int
	ch   chan Task
	done chan struct{}
	db   *repo.DB
	ctx  context.Context
	mu   *sync.Mutex
}

func NewInputWorker(ch chan Task, done chan struct{}, ctx context.Context, mu *sync.Mutex) *InputWorker {
	index := 0
	ticker := time.NewTicker(10 * time.Second)
	return &InputWorker{
		ch:     ch,
		done:   done,
		index:  index,
		ticker: ticker,
		ctx:    ctx,
		mu:     mu,
	}
}

func NewOutputWorker(id int, ch chan Task, done chan struct{}, ctx context.Context, db *repo.DB, mu *sync.Mutex) *OutputWorker {
	return &OutputWorker{
		id:   id,
		ch:   ch,
		done: done,
		ctx:  ctx,
		db:   db,
		mu:   mu,
	}
}

func (w *InputWorker) Do(t Task) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.ch <- t
	w.index++
	log.Println(w.index)
	if w.index == 20 {
		w.done <- struct{}{}
		w.index = 0
	}
}

func (w *InputWorker) Loop() error {
	for {
		select {
		case <-w.ctx.Done():
			w.ticker.Stop()
			return nil
		case <-w.ticker.C:
			w.mu.Lock()
			w.done <- struct{}{}
			w.index = 0
			w.mu.Unlock()
		}
	}
}

func (w *OutputWorker) Do() error {
	models := make([]entity.ModelURLForDelete, 0, 20)
	var URLForDelete entity.ModelURLForDelete

	for {
		select {
		case <-w.ctx.Done():
			return nil
		case <-w.done:
			if len(w.ch) == 0 {
				break
			}
			for task := range w.ch {
				for _, url := range task.IDsForDelete {
					URLForDelete.UserID = task.UserID
					URLForDelete.ShortURL = url
					models = append(models, URLForDelete)
				}
				if len(w.ch) == 0 {
					if err := w.db.BatchDeleteShortURLs(models); err != nil {
						return err
					}
					models = nil
					break
				}
			}
		}
	}
}
