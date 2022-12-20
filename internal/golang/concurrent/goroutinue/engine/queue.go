package main

// QueuedScheduler :Scheduler with requestChan and workerChan
type QueuedScheduler struct {
	requestChan chan string
	workerChan  chan chan string
}

func (s *QueuedScheduler) Submit(r string) {
	s.requestChan <- r
}

func (s *QueuedScheduler) WorkerChan() chan string {
	return make(chan string)
}

func (s *QueuedScheduler) WorkerReady(w chan string) {
	s.workerChan <- w
}

func (s *QueuedScheduler) Run() {
	// init requestChan and workerChan
	s.requestChan = make(chan string)
	s.workerChan = make(chan chan string)

	go func() {
		// create queue request and worker
		var requestQ []string
		var workerQ []chan string
		for {
			var activeRequest string
			var activeWorker chan string
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case r := <-s.requestChan:
				// 收到 request 加入队列进行排队
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				// 收到 worker 加入队列进行排队
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				// 若 activeRequest and activeWorker 都有数据,将请求发送至工作队列
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
