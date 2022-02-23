package scheduler

import "study/internal/crawler/engine"

// QueuedScheduler :Scheduler with requestChan and workerChan
type QueuedScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request
}

func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

func (s *QueuedScheduler) Run() {
	// init requestChan and workerChan
	s.requestChan = make(chan engine.Request)
	s.workerChan = make(chan chan engine.Request)

	go func() {
		// create queue request and worker
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
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
