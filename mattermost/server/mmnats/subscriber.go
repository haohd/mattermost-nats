package mmnats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	MaxSubscriberMessages = 512
)

type Subscriber interface {
	Subscribe(subject, group string, handler MessageHandler[*Message], autoReply bool) error
	Unsubscribe()
}

func NewSubscriber(con Connection) Subscriber {
	s := &subscriber{
		con: con,
		ch:  make(chan *Message, MaxSubscriberMessages),
	}
	return s
}

type subscriber struct {
	con Connection
	ch  chan *Message
	sub *nats.Subscription
}

func (s *subscriber) Subscribe(subject, group string, handler MessageHandler[*Message], autoReply bool) (err error) {
	if s.con == nil {
		log.Fatal("NATS has not been connected!")
	}
	subject = s.con.SubjectPrefix(subject)
	group = s.con.NamePrefix(group)
	s.sub, err = s.con.GetConnection().ChanQueueSubscribe(subject, group, s.ch)
	if err != nil {
		log.Printf("Failed to subscribe subject [%s], group [%s]", subject, group)
		return
	}

	log.Printf("Subscribed event for subject [%s], group [%s]", subject, group)

	go func() {
		for msg := range s.ch {
			if msg != nil {
				go func() {
					err := handler(msg)
					if autoReply {
						if err == nil {
							err := msg.Ack()
							if err != nil {
								log.Printf("Failed to ACK subscribe subject [%s], error [%v]", subject, err)
							}
						} else {
							err := msg.Nak()
							if err != nil {
								log.Printf("Failed to NAK subscribe subject [%s], error [%v]", subject, err)
							}
						}
					} else {
						if err != nil {
							log.Printf("Failed to subscribe subject [%s], error [%v]", subject, err)
						}
					}
				}()
			}
		}
		log.Printf("Unsubscribed event for subject [%s], group [%s]", subject, group)
	}()
	return err
}

func (s *subscriber) Unsubscribe() {
	if s.sub == nil {
		return
	}
	err := s.sub.Drain()
	if err != nil {
		log.Printf("Failed to drain subscribe connection subject [%s], error [%v]", s.sub.Subject, err)
	}
	if s.ch != nil {
		// Wait to consume all messages in the message channel before exit consumer routine
		for {
			select {
			case status := <-s.sub.StatusChanged():
				if status == nats.SubscriptionClosed {
					close(s.ch)
					return
				}
			case <-time.After(time.Second * 5):
				log.Printf("Unsubscribe timeout on subject: %v", s.sub.Subject)
				return
			}
		}
	}
}
