package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

type DataChannel chan DataEvent

type DataChannelSlice []DataChannel

type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chs, found := eb.subscribers[topic]; found {
		channels := append(DataChannelSlice{}, chs...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

func (eb *EventBus) Listen(interrupt chan os.Signal, ch1, ch2, ch3 chan DataEvent) {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		outer: for {
			select {
			case d := <-ch1:
				go printDataEvent("ch1", d)
			case d := <-ch2:
				go printDataEvent("ch2", d)
			case d := <-ch3:
				go printDataEvent("ch3", d)
			case <-interrupt:
				fmt.Printf("\nShoutdown!\n")
				break outer
			}
		}
		wg.Done()
	}()

	wg.Wait()
	os.Exit(0)
}

var eb = &EventBus{
	subscribers: map[string]DataChannelSlice{},
}

func printDataEvent(ch string, data DataEvent) {
	fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}

func publishTo(topic string, data string) {
	for {
		eb.Publish(topic, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	ch1 := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	ch3 := make(chan DataEvent)
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	eb.Subscribe("topic1", ch1)
	eb.Subscribe("topic2", ch2)
	eb.Subscribe("topic2", ch3)

	go publishTo("topic1", "Hello I am Boris from topic 1")
	go publishTo("topic2", "Hello mazafaka from topic 2")

	eb.Listen(interrupt, ch1, ch2, ch3)

}
