package main

type Dispatcher struct {
	maxWorkers int
	WorkerPool chan chan string
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{maxWorkers: maxWorkers, WorkerPool: make(chan chan string, maxWorkers)}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go func() {
		for {
			select {
			case job := <-JobQueue:
				go func(job string) {
					jobChannel := <-d.WorkerPool // 尝试获得可得到的 worker job channel.this will block until a worker is idle
					jobChannel <- job            // dispatch job to job channel
				}(job)
			}
		}
	}()
}
