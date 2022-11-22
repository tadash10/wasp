package publisher

import (
	"github.com/iotaledger/hive.go/core/events"
)

type BaseEvent struct{}

type ChainEvent struct {
	MessageType string
	ChainID     string
	Message     interface{}
}

func NewChainEvent(messageType, chainID string, message interface{}) *ChainEvent {
	return &ChainEvent{
		MessageType: messageType,
		ChainID:     chainID,
		Message:     message,
	}
}

var Event = events.NewEvent(func(handler interface{}, msg ...interface{}) {
	callback := handler.(func(event *ChainEvent))
	callback(msg[0].(*ChainEvent))
})

func Publish(event *ChainEvent) {
	Event.Trigger(event)
}
