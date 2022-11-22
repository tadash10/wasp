package eventmessages

import (
	"github.com/iotaledger/wasp/packages/publisher"
)

type RequestOut struct {
	RequestID  string
	RequestIDs int
	StateIndex uint32
}

func NewChainEventRequestOut(chainID string, message RequestOut) *publisher.ChainEvent {
	return publisher.NewChainEvent("request_out", chainID, message)
}

type StateUpdate struct {
	RequestIDs    int
	StateHash     string
	StateIndex    uint32
	StateOutputID string
}

func NewChainEventStateUpdate(chainID string, message StateUpdate) *publisher.ChainEvent {
	return publisher.NewChainEvent("state", chainID, message)
}

type Rotation struct {
	RequestIDs    int
	StateHash     string
	StateIndex    uint32
	StateOutputID string
}

func NewChainEventRotation(chainID string, message Rotation) *publisher.ChainEvent {
	return publisher.NewChainEvent("rotate", chainID, message)
}

type VMMessage struct {
	EventMessage string
}

func NewChainEventVMMessage(chainID string, message VMMessage) *publisher.ChainEvent {
	return publisher.NewChainEvent("vmmsg", chainID, message)
}

func NewChainEventDismissed(chainID string) *publisher.ChainEvent {
	return publisher.NewChainEvent("rotate", chainID, nil)
}
