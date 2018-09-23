package worker

import (
	"sync"
	"time"
)

type Process func(interface{})

type AutoScalingWorker struct {
	MinWorker     int
	MaxWorker     int
	QueueDepth    int
	CurrentWorker int
	Process       Process
	Queue         chan interface{}
	Interval      time.Duration

	stop       chan int
	workerStop chan int
	lock       sync.RWMutex
}

func (auto *AutoScalingWorker) Start() {
	auto.workerStop = make(chan int)
	auto.stop = make(chan int)

	t := time.NewTicker(auto.Interval)
	for {
		select {
		case <-t.C:
			if len(auto.Queue) > auto.QueueDepth {
				auto.Expansion()
			} else {
				auto.Shrinkage()
			}
		case <-auto.stop:
			return
		}
	}

}

func (auto *AutoScalingWorker) Expansion() {
	if auto.CurrentWorker < auto.MaxWorker {
		go auto.job()
	}
}
func (auto *AutoScalingWorker) Shrinkage() {
	if auto.CurrentWorker > auto.MinWorker {
		auto.workerStop <- 1

	}
}

func (auto *AutoScalingWorker) Stop() {
	close(auto.workerStop)
	close(auto.stop)
}

func (auto *AutoScalingWorker) job() {
	auto.lock.Lock()
	auto.CurrentWorker = auto.CurrentWorker + 1
	auto.lock.Unlock()
	for {
		select {
		case <-auto.workerStop:
			goto close

		case data, ok := <-auto.Queue:
			if !ok {
				goto close
			}
			auto.Process(data)
		}
	}

close:
	auto.lock.Lock()
	auto.CurrentWorker = auto.CurrentWorker - 1
	auto.lock.Unlock()
}
