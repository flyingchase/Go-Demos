package goroutinepooldemo

import (
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	capacity int32
	// current running goroutines
	running int32
	// set the expired time for a worker
	expiryDuration time.Duration
	workers        []*Worker
	// flag to notice to close pool itself
	release int32
	lock    sync.Mutex
	// ensure releaing  to do once
	once sync.Once

	// speeds up the obtainment of the usable work in func retrieveWorker
	workerCache sync.Pool
	// cond for waiting idleWorker
	cond sync.Cond
}

func NewTimingPool(size int, expiry int) (*Pool, error) {
	if size <= 0 {
		return nil, ErrInvalidPoolSize
	}
	if expiry <= 0 {
		return nil, ErrInvalidPoolExpiry
	}
	p := &Pool{
		capacity:       int32(size),
		expiryDuration: time.Duration(expiry) * time.Second,
	}
	return p, nil
}

func NewPool(size int) (*Pool, error) {
	return NewTimingPool(size, DefaultCleanIntervalTime)
}

// submits a task to this pool
func (p *Pool) Submit(task func()) error {
	// 判断该 pool 是关闭
	if Closed == atomic.LoadInt32(&p.release) {
		return ErrPoolClosed
	}
	// 检索出可用的 worker 后将 task 绑定到该 worker的 task chan 上
	p.retrieveWorker().task <- task
	return nil
}

// retrieveWorker 检索可用的 worker to run the task
func (p *Pool) retrieveWorker() *Worker {
	var w *Worker
	p.lock.Lock()
	idleWorkers := p.workers
	n := len(idleWorkers) - 1
	// 闲置 worker 队列中存在
	if n >= 0 {
		w = idleWorkers[n]
		// 取出闲置 worker 后将原p.workers 对应的 worker 置空
		idleWorkers[n] = nil
		p.workers = idleWorkers[:n]
		p.lock.Unlock()
	} else if p.Runing() < p.Cap() {
		// workers 切片中无空闲worker但辅助池中存在
		p.lock.Lock()
		// 从 workerCache辅助池中随机获得可用的 worker
		if cacheWork := p.workerCache.Get(); cacheWork != nil {
			w = cacheWork.(*Worker)
		} else {
			w = &Worker{
				pool: p,
				task: make(chan func(), 1),
			}
		}
		w.run()
	} else {
		//阻塞判断
		for {
			p.cond.Wait()
			l := len(p.workers) - 1
			// 出现空闲 idleWorker则取出队尾
			if l < 0 {
				continue
			}

			w = p.workers[l]
			p.workers[l] = nil
			p.workers = p.workers[:l]
			break
		}
		p.lock.Unlock()
	}
	return w
}

// put worker into pool to recycle goroutines
func (p *Pool) revertWorker(worker *Worker) bool {
	// 先检测 pool 是否关闭
	if Closed == atomic.LoadInt32(&p.release) {
		return false
	}
	// 入池时更新 worker 的 time
	worker.recycleTime = time.Now()
	p.lock.Lock()
	p.workers = append(p.workers, worker)

	// retrieveWorker() stuck 时是存在可用 worker
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

// clear expired workers periodically
func (p *Pool) periodicallyPurge() {
	// 设定过期时间
	heartBeat := time.NewTicker(p.expiryDuration)
	defer heartBeat.Stop()

	for range heartBeat.C {
		if Closed == atomic.LoadInt32(&p.release) {
			break
		}
		currentTime := time.Now()
		p.lock.Lock()
		idleWorkers := p.workers
		n := -1
		for i, w := range idleWorkers {
			if currentTime.Sub(w.recycleTime) <= p.expiryDuration {
				break
			}
			n = i
			w.task <- nil
			idleWorkers[i] = nil
		}
		if n > -1 {
			if n >= len(idleWorkers)-1 {
				p.workers = idleWorkers[:0]
			} else {
				p.workers = idleWorkers[n+1:]
			}
		}
		p.lock.Unlock()
	}

}

// package cap run free and close
func (p *Pool) Runing() int {
	return int(atomic.LoadInt32(&p.running))
}
func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.capacity))
}

// return available goroutines to work
func (p *Pool) Free() int {
	return int(atomic.LoadInt32(&p.capacity)) - int(atomic.LoadInt32(&p.running))
}
func (p *Pool) Tune(size int) {
	if size == p.Cap() {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
	// 将超出目新 size 的 worker 全部置空
	diff := p.Runing() - size
	for i := 0; i < diff; i++ {
		p.retrieveWorker().task <- nil
	}
}

// release close this pool
func (p *Pool) Release() error {
	// once
	p.once.Do(func() {
		atomic.StoreInt32(&p.release, int32(1))
		p.lock.Lock()
		idleWorkers := p.workers
		for i, w := range idleWorkers {
			w.task <- nil
			idleWorkers[i] = nil
		}
		p.workers = nil
		p.lock.Unlock()
	})
	return nil

}

// increase the number of current running goroutine
func (p *Pool) incRunning() {
	atomic.AddInt32(&p.running, int32(1))
}

// decrease the number of current running goroutines
func (p *Pool) decRunning() {
	atomic.AddInt32(&p.running, int32(-1))
}
