package mmnats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Message = nats.Msg // For generic subscriber
type PushMessage = *nats.Msg
type PullMessage = jetstream.Msg
type MessageHandler[T any] func(T) error

type Consumer[T any] interface {
	Consume(subjects []string, handler MessageHandler[T]) error
	Unsubscribe()
}
