package event

import (
	"github.com/asaskevich/EventBus"
	"github.com/google/wire"
)

// ProviderSet is EventBus providers.
var ProviderSet = wire.NewSet(NewEventBus, NewBusPublisher)

var (
	eventBus EventBus.Bus
)

func init() {
	eventBus = EventBus.New()
}

func NewEventBus() EventBus.Bus {
	return eventBus
}

func NewBusPublisher() EventBus.BusPublisher {
	return eventBus
}
