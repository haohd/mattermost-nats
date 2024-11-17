package mmnats

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
)

var (
	natsConn *natsConnection
	once     sync.Once
)

type Connection interface {
	Close()
	GetConnection() *nats.Conn
	NamePrefix(name string) string
	SubjectPrefix(name string) string
}

type natsConnection struct {
	cfg        *NatsConfig
	connection *nats.Conn
}

func NatsConnection() Connection {
	once.Do(func() {
		cfg := &NatsConfig{}
		lo.Must0(envconfig.Process("", cfg))

		if cfg.URL == "" {
			// Do nothing if NATS is not enabled
			return
		}

		opts := nats.GetDefaultOptions()
		opts.ClosedCB = func(c *nats.Conn) {
			// If still lost connection after retries, just panic app HERE
			log.Fatal("NATS connection closed.")
		}
		opts.Servers = []string{cfg.URL}

		var connection *nats.Conn
		_, err := lo.Attempt(10, func(index int) error {
			var err error
			connection, err = opts.Connect()
			if err != nil {
				log.Printf("Failed to connect to NATS: %v", err)
				time.Sleep(time.Second)
				return err
			}
			return nil
		})
		lo.Must0(err)

		log.Println("Connected to NATS successfully.")

		natsConn = &natsConnection{
			cfg:        cfg,
			connection: connection,
		}
	})
	return natsConn
}

func (c *natsConnection) Close() {
	if c.connection != nil {
		err := c.connection.Drain()
		if err != nil {
			log.Printf("Drain NATS connection has error: %s", err)
		}
		c.connection.Close()
	}
}

func (c *natsConnection) GetConnection() *nats.Conn {
	return c.connection
}

func (c *natsConnection) NamePrefix(name string) string {
	return fmt.Sprintf("%s:%s", c.cfg.EnvPrefix, name)
}

func (c *natsConnection) SubjectPrefix(name string) string {
	return fmt.Sprintf("%s.%s", c.cfg.EnvPrefix, name)
}
