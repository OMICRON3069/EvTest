package main

import (
	"EvTest/evBus"
	"fmt"
	"math/rand"
	"time"
)

//TODO: flag parse, working directory initiate,

func main() {
	theFirstBus := evBus.New()

	ch1 := make(chan evBus.Event)
	ch2 := make(chan evBus.Event)
	ch3 := make(chan evBus.Event)

	theFirstBus.SubScribe("http:sb", ch1, printDataEvent)
	theFirstBus.SubScribe("http:dsb", ch2, printDataEvent)
	theFirstBus.SubScribe("http:dsb", ch3, printDataEvent)
	fmt.Println("??")
	go func(topic string, data string) {
		fmt.Println("start this")
		for {
			_ = theFirstBus.Publish(topic, data)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}("http:sb", "Hi, Topic1")

	go func(topic string, data string) {
		for {
			_ = theFirstBus.Publish(topic, data)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}("http:dsb", "sucker topic2")
	time.Sleep(time.Second * 5)
	fmt.Println("??x")
}

func printDataEvent(ch evBus.EventChannel) {
	data := <-ch
	fmt.Printf("Topic: %s; DataEvent: %v\n", data.Topic, data.Data)
}
