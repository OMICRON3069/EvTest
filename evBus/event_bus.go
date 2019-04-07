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
	ID    uint64
}

type EventChannel chan Event

type Bus struct {
	subscriber map[string][]EventChannel
	sync.RWMutex
}

type SubQueue struct {
	//TODO: loop this queue to receive callback
	//maybe that's not the most effective solution
	//see https://stackoverflow.com/questions/3398490/checking-if-a-channel-has-a-ready-to-read-value-using-go
	//and https://stackoverflow.com/questions/19992334/how-to-listen-to-n-channels-dynamic-select-statement
}

func (b *Bus) SubScribe(topic string, ch EventChannel) {
	b.Lock()
	if pre, found := b.subscriber[topic]; found {
		b.subscriber[topic] = append(pre, ch)
	} else {
		b.subscriber[topic] = append([]EventChannel{}, ch)
	}
	b.Unlock()
}

//This method will publish data to specific topic
//and it will return error if data interface is not data type
func (b *Bus) Publish(topic string, data interface{}) error {
	b.RLock()
	defer b.RUnlock()

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
		if ch, found := b.subscriber[topic]; found {
			//create a new slice here to be passed to func
			chs := append([]EventChannel{}, ch...)
			//use a separate routine to callback
			go func(data Event, chs []EventChannel) {
				for _, ch := range chs {
					ch <- data
				}
			}(Event{Data: data, Topic: topic}, chs)
		}
	} //TODO: what if there is no subscriber?
	return nil
}

func New() *Bus {
	return &Bus{
		subscriber: make(map[string][]EventChannel),
	}
}
