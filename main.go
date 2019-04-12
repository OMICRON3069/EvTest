package main

import (
	"EvTest/evBus"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//TODO: flag parse, working directory initiate,

func main() {
	theFirstBus := evBus.New()

	var wg sync.WaitGroup
	ch1 := make(chan evBus.Event)
	ch2 := make(chan evBus.Event)
	ch3 := make(chan evBus.Event)

	theFirstBus.SubScribe("http:sb", ch1, printDataEvent)
	theFirstBus.SubScribe("http:dsb", ch2, printDataEvent)
	theFirstBus.SubScribe("http:dsb", ch3, printDataEvent)
	wg.Add(1)
	go func(topic string, data string) {

		fmt.Println("start 1")
		for {
			_ = theFirstBus.Publish(topic, data)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		//wg.Done()
	}("http:sb", "Hi, Topic1")

	wg.Add(1)
	go func(topic string, data string) {

		fmt.Println("start 2")
		for {
			_ = theFirstBus.Publish(topic, data)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		//wg.Done()
	}("http:dsb", "sucker topic2")
	wg.Wait()
	fmt.Println("??x")
}

func printDataEvent(ch evBus.EventChannel) {
	data := <-ch
	fmt.Printf("Topic: %s; DataEvent: %v\n", data.Topic, data.Data)
}
