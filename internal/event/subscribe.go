package event

import (
	"github.com/asaskevich/EventBus"
)

type listener struct {
	topic string
	fn    interface{}
}

var listeners = []listener{
	//	listener{topic: "", fn: nil},
}

func LoadListener(bus EventBus.Bus) (err error) {
	for _, listener := range listeners {
		err = bus.Subscribe(listener.topic, listener.fn)
		if err != nil {
			return
		}
	}
	return
}
