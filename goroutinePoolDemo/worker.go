package goroutinepooldemo

import (
	"time"
)

type Worker struct {

	// 代表该worker所属的池
	pool *Pool
	task chan func()
	// 当该 worker 入队列时会更新
	// 入队时间
	recycleTime time.Time
}

func (w *Worker) run() {
	w.pool.incRunning()
	go func() {

		defer func() {
			if p := recover(); p != nil {
				w.pool.decRunning()
				w.pool.workerCache.Put(w)
				// panic deal

			}
		}()

		for f := range w.task {
			if f == nil {
				w.pool.decRunning()
				w.pool.workerCache.Put(w)
				return
			}
			f()
			if ok := w.pool.revertWorker(w); !ok {
				break
			}

		}

	}()
}
