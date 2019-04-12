package evBus

//This is the test project

import (
	"EvTest/jankyError"
	"reflect"
	"sync"
)

type Event struct {
	Data  interface{}
	Topic string
}

type EventChannel chan Event

type Bus struct {
	subscriber map[string][]SubQueue
	sync.RWMutex
}

type SubQueue struct {
	Holder func(ch EventChannel)
	Messenger EventChannel
}


func (bus *Bus) SubScribe(topic string, ch EventChannel, holder func(ch EventChannel)) {
	bus.Lock()
	defer bus.Unlock()
	if pre, found := bus.subscriber[topic]; found {
		bus.subscriber[topic] = append(pre,SubQueue{holder,ch})
	} else {
		bus.subscriber[topic] = append([]SubQueue{}, SubQueue{holder,ch})
	}
}

func (bus *Bus)UnSubscribe(topic string, ch EventChannel, holder func())  {
	//TODO
}

//This method will publish data to specific topic
//and it will return error if data interface is not data type
func (bus *Bus) Publish(topic string, data interface{}) error {
	bus.RLock()
	defer bus.RUnlock()

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
			contacts := append([]SubQueue{}, contact...)
			//use a separate routine to callback
			go func(data Event, contacts []SubQueue) {
				for _, contact := range contacts {
					//TODO:check if channel closed
					go contact.Holder(contact.Messenger)
					contact.Messenger <- data
				}
			}(Event{Data: data, Topic: topic}, contacts)
		}
	} //TODO: what if there is no subscriber?
	return nil
}

//maybe I need a better way to create a new bus
func New() *Bus {
	return &Bus{
		subscriber: make(map[string][]SubQueue),
	}
}
