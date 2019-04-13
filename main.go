package main

import (
	"EvTest/evBus"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//TODO: flag parse, working directory initiate,

func main() {
	theFirstBus := evBus.New()

	var (
		wg sync.WaitGroup
	)

	theFirstBus.Subscribe("http:sb", printDataEvent)
	theFirstBus.Subscribe("http:dsb",printDataEvent)
	theFirstBus.Subscribe("http:dsb", printDataEvent)

	wg.Add(1)
	go func(topic string, data string) {
		log.Println("start 1")
		for {
			theFirstBus.Publish(topic, data)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		//wg.Done()
	}("http:sb", "Hi, Topic1")

	wg.Add(1)
	go func(topic string, data string) {
		log.Println("start 2")
		for {
			theFirstBus.Publish(topic, data)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		//wg.Done()
	}("http:dsb", "sucker topic2")
	wg.Wait()
}

func printDataEvent(data string) {
	fmt.Println(data)
}
