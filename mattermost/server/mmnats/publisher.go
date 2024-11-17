package mmnats

import (
	"time"
)

type Publisher interface {
	Publish(subject string, data []byte) error
	Request(subject string, data []byte, timeout time.Duration) (*Message, error)
}

func NewPublisher(con Connection) Publisher {
	return &natsPublisher{
		connection: con,
	}
}

type natsPublisher struct {
	connection Connection
}

func (n *natsPublisher) Publish(subject string, data []byte) error {
	return n.connection.GetConnection().Publish(n.connection.SubjectPrefix(subject), data)
}

func (n *natsPublisher) Request(subject string, data []byte, timeout time.Duration) (*Message, error) {
	return n.connection.GetConnection().Request(n.connection.SubjectPrefix(subject), data, timeout)
}
