package crawler

type Queue struct {
	jobChan  chan Request
	jobChanQ chan chan Request
}

func NewQueuedScheduler() *Queue {
	return &Queue{jobChan: make(chan Request), jobChanQ: make(chan chan Request)}
}

func (s *Queue) Run() {
	var jobs []Request
	var jobChans []chan Request
	for {
		var activeJob Request
		var activeJobChan chan Request
		if len(jobs) > 0 && len(jobChans) > 0 {
			activeJob = jobs[0]
			activeJobChan = jobChans[0]
		}

		select {
		case job := <-s.jobChan:
			jobs = append(jobs, job)
		case jobChan := <-s.jobChanQ:
			jobChans = append(jobChans, jobChan)
		case activeJobChan <- activeJob:
			jobs = jobs[1:]
			jobChans = jobChans[1:]
		}
	}
}

func (s *Queue) Submit(job Request) {
	s.jobChan <- job
}

func (s *Queue) JobChan() chan Request {
	return make(chan Request)
}

func (s *Queue) Dispatcher(jobChan chan Request) {
	s.jobChanQ <- jobChan
}
