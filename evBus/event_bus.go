package evBus

//This is the test project

import (
	"EvTest/jankyError"
	"log"
	"reflect"
	"sync"
	"time"
)

type Event struct {
	Data  interface{}
	Topic string
}

type EventChannel chan Event

type Bus struct {
	subscriber   map[string][]EventChannel
	queue        []SubQueue
	self, quLock sync.RWMutex
}

type SubQueue struct {
	Holder    func(ch reflect.Value)
	Messenger EventChannel
}

func (bus *Bus) StartLoop() {
	bus.quLock.RLock()
	cases := make([]reflect.SelectCase, len(bus.queue))
	for i, ch := range bus.queue {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.Messenger)}
	}
	bus.quLock.RUnlock()

	go bus.ticking()

	for {
		bus.quLock.RLock()
		if chosen, value, ok := reflect.Select(cases); ok {
			go bus.queue[chosen].Holder(value)
		}
		bus.quLock.RLock()
	}
}

func (bus *Bus) ticking() {
	click := make(chan Event)
	gugu := func(ch reflect.Value) {
		data := ch.Elem()
		log.Println(data.FieldByName("Data"))
	}
	bus.SubScribe("bus:ticking", click, gugu)

	for {
		time.Sleep(time.Second * 2)
		_ = bus.Publish("bus:ticking", "it's time")
	}
}

func (bus *Bus) SubScribe(topic string, ch EventChannel, holder func(ch reflect.Value)) {
	bus.self.Lock()
	defer bus.self.Unlock()

	if pre, found := bus.subscriber[topic]; found {
		bus.subscriber[topic] = append(pre, ch)
	} else {
		bus.subscriber[topic] = append([]EventChannel{}, ch)
	}
	//start make queue
	bus.quLock.Lock()
	bus.queue = append([]SubQueue{}, SubQueue{
		Holder:    holder,
		Messenger: ch,
	})
	bus.quLock.Unlock()
}

func (bus *Bus) UnSubscribe(topic string, ch EventChannel) {
	//TODO
}

//This method will publish data to specific topic
//and it will return error if data interface is not data type
func (bus *Bus) Publish(topic string, data interface{}) error {
	bus.self.RLock()
	defer bus.self.RUnlock()

	//checking type of data
	if reflect.TypeOf(data).Kind() == reflect.Func {
		//TODO: figure out does this "&" is needed
		return &jankyError.TheError{
			Code:    jankyError.NotDataCode,
			Message: jankyError.NotData,
			Detail:  nil,
		}
	} else {
		//if this topic has subscriber
		if contact, found := bus.subscriber[topic]; found {
			//create a new slice here to be passed to func
			contacts := append([]EventChannel{}, contact...)
			//use a separate routine to callback
			go func(data Event, contacts []EventChannel) {
				for _, contact := range contacts {
					//TODO:check if channel closed
					contact <- data
				}
			}(Event{Data: data, Topic: topic}, contacts)
		}
	} //TODO: what if there is no subscriber?
	return nil
}

//maybe I need a better way to create a new bus
func New() *Bus {
	return &Bus{
		subscriber: make(map[string][]EventChannel),
		queue:      make([]SubQueue, 0),
	}
}
