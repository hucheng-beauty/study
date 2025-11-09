package main

import "time"

type Session struct{ MsgChan chan string }

func (s *Session) Response(err error) {}

type Worker struct {
    sess *Session
    done bool
}

func (w *Worker) handle(message string) error { return nil }

func (w *Worker) clean() {}

func (w *Worker) checkTimeout() error { return nil }

func (w *Worker) Start() {
    tick := time.NewTicker(time.Second)
    var err error
    for !w.done && err == nil {
        select {
        case msg, ok := <-w.sess.MsgChan:
            if ok {
                err = w.handle(msg)
            }
        case <-tick.C:
            err = w.checkTimeout()
        }
    }

    if err != nil {
        w.sess.Response(err)
    }
    w.clean()
}

/*
	整体数据流向:
	request  ==>  scheduler ==>  session  ==> worker <==> topology
												||
											   task(sync、async、max-10) ==> taskPool(100)
								PreProcess ==> task ==> PostProcess ==> 收集结果 ==> response
*/
func main() {

    limitSize := 100
    workerSize := 10000

    workerBuffer := make(chan Worker, workerSize)
    limitChan := make(chan struct{}, limitSize)

    go func() {
        for {
            select {
            case worker, ok := <-workerBuffer:
                if ok {
                    limitChan <- struct{}{}
                    go func(w Worker) {
                        defer func() { <-limitChan }()
                        w.Start()
                    }(worker)
                }
            }
        }
    }()
}
