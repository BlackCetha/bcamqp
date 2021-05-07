package bcamqp

import "github.com/streadway/amqp"

// ExchangeType specifies the type of the exchange to be created
type ExchangeType string

// Lists the different exchange types
const (
	Direct  ExchangeType = "direct"
	Fanout  ExchangeType = "fanout"
	Topic   ExchangeType = "topic"
	Headers ExchangeType = "headers"
)

type BrokerOptions struct {
	Encrypted     bool // Wether to use AMQPs
	Address       string
	User          string
	Password      string
	AutoTimestamp bool // Wether to set the (unset) timestamp field when publishing messages
}

// QueueOptions holds options for queue creation
type QueueOptions struct {
	Name             string
	Durable          bool
	Exclusive        bool
	ConsumeExclusive bool
	ConsumerName     string
}

// ExchangeOptions holds options for exchange creation
type ExchangeOptions struct {
	Name    string
	Type    ExchangeType
	Durable bool
}

type propagationHeaders amqp.Table

func (p propagationHeaders) Get(key string) string {
	v, _ := p[key]
	vstr, _ := v.(string)
	return vstr // use default value
}

func (p propagationHeaders) Set(key, value string) {
	p[key] = value
}

func (p propagationHeaders) Keys() []string {
	out := make([]string, 0, len(p))

	for k := range p {
		out = append(out, k)
	}

	return out
}
