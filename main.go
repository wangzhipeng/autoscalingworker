package main

import (
	"fmt"
	"time"

	"github.com/wangzhipeng/autoscalingworker/worker"
)

func main() {
	c := make(chan interface{}, 1000)

	auto := worker.AutoScalingWorker{
		MinWorker:  0,
		MaxWorker:  10,
		QueueDepth: 8,
		Interval:   time.Second,
		Process: func(i interface{}) {
			time.Sleep(time.Millisecond * 100)
		},
		Queue: c,
	}

	go auto.Start()
	go func() {
		for n := 0; n < 5; n++ {
			for i := 0; i < 30; i++ {
				c <- i
			}
			time.Sleep(time.Second)
		}

	}()
	go func() {
		time.Sleep(time.Second * 15)
		for n := 0; n < 5; n++ {
			for i := 0; i < 30; i++ {
				c <- i
			}

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
