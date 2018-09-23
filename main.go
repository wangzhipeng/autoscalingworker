package main

import (
	"fmt"
	"time"

	"github.com/wangzhipeng/autoscalingworker/worker"
)

func main() {
	c := make(chan interface{}, 100000)

	auto := worker.AutoScalingWorker{
		MinWorker:  1,
		MaxWorker:  10,
		QueueDepth: 100,
		Interval:   time.Second,
		Process: func(i interface{}) {
			time.Sleep(time.Millisecond * 100)
		},
		Queue: c,
	}

	go auto.Start()
	go func() {
		for n := 0; n < 10; n++ {
			for i := 0; i < 100; i++ {
				c <- i
			}
			time.Sleep(time.Second)
		}

		for n := 0; n < 10; n++ {
			for i := 0; i < 50; i++ {
				c <- i
			}
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			fmt.Println("auto worker:", auto.CurrentWorker, " QueueDepth:", len(c))
			time.Sleep(time.Second)
		}

	}()

	time.Sleep(time.Minute)
}
