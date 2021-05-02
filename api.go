package bcamqp

import (
	"fmt"
	"net/url"

	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/blackcetha/bcamqp/adaptor/streadway"
)

var adaptorInstance adaptor.Server = streadway.Streadway{}

func Open(options BrokerOptions, onClose func(err error)) (*Broker, error) {
	amqpURL := url.URL{
		Scheme: "amqp",
		Host:   options.Address,
		User:   url.UserPassword(options.User, options.Password),
	}
	if options.Encrypted {
		amqpURL.Scheme = "amqps"
	}

	var b Broker
	err := b.open(amqpURL.String(), options.AutoTimestamp, onClose)
	if err != nil {
		return nil, fmt.Errorf("connect to broker: %w", err)
	}

	return &b, nil
}
