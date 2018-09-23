# autoscalingworker

## 介绍
传统基于channel的生产者与消费者模式，采用固定消费者个数形式，在面对生产数量猛增情况时，不能做到动态伸缩，容易致使队列积压。基于此 autoscalingworker 采用检测队列积压情况动态创建与销毁worker，来保证消费与避免性能的浪费。

## demo
``` go
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
```
输出结果：
```
auto worker: 0  QueueDepth: 60
auto worker: 2  QueueDepth: 79
auto worker: 3  QueueDepth: 89
auto worker: 4  QueueDepth: 89
auto worker: 5  QueueDepth: 49
auto worker: 5  QueueDepth: 0
auto worker: 3  QueueDepth: 0
auto worker: 2  QueueDepth: 0
auto worker: 1  QueueDepth: 0
auto worker: 0  QueueDepth: 150
auto worker: 1  QueueDepth: 149
auto worker: 2  QueueDepth: 138
auto worker: 3  QueueDepth: 118
auto worker: 4  QueueDepth: 88
auto worker: 5  QueueDepth: 48
auto worker: 4  QueueDepth: 0
auto worker: 3  QueueDepth: 0
auto worker: 2  QueueDepth: 0
auto worker: 1  QueueDepth: 0
auto worker: 0  QueueDepth: 0
```