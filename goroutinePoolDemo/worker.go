package goroutinepooldemo

import "time"

type Worker struct {

	// 代表该worker所属的池
	pool *Pool
	task chan func()
	// 当该 worker 入队列时会更新
	// 入队时间
	recycleTime time.Time
}

func (w *Worker) run() {

}
