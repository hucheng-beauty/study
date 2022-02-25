package message_exchange

// Scheduler : Scheduler with requestChan and workerChan
type Scheduler struct {
	requestChan chan interface{}
	workerChan  chan chan interface{}
}

func (s *Scheduler) Submit(r interface{}) {
	s.requestChan <- r
}

func (s *Scheduler) WorkerChan() chan interface{} {
	return make(chan interface{})
}

func (s *Scheduler) WorkerReady(w chan interface{}) {
	s.workerChan <- w
}

func (s *Scheduler) Run() {
	// init requestChan and workerChan
	s.requestChan = make(chan interface{})
	s.workerChan = make(chan chan interface{})

	go func() {
		// create queue request and worker
		var requestQ []interface{}
		var workerQ []chan interface{}
		for {
			var activeRequest interface{}
			var activeWorker chan interface{}
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
