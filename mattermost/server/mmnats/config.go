package mmnats

type NatsConfig struct {
	URL       string `envconfig:"NATS_URL" default:""`
	EnvPrefix string `envconfig:"NATS_PREFIX" default:"mm"`
}
