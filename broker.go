package bcamqp

import (
	"fmt"

	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/streadway/amqp"
)

type Broker struct {
	connection    adaptor.Connection
	mainChan      adaptor.Channel
	autoTimestamp bool
}

func (b *Broker) open(url string, autoTimestamp bool, onClose func(err error)) error {
	var err error
	b.connection, err = adaptorInstance.Dial(url)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	b.mainChan, err = b.connection.Channel()
	if err != nil {
		return fmt.Errorf("open management channel: %w", err)
	}

	b.autoTimestamp = autoTimestamp

	// convert channel read to callback for easier external API
	disconnectChan := make(chan *amqp.Error)
	b.connection.NotifyClose(disconnectChan)
	go func() {
		onClose(<-disconnectChan)
	}()

	return nil
}

func (b *Broker) Exchange(options ExchangeOptions) (*Exchange, error) {
	// we need a special case for the default exchange since RabbitMQ
	// won't let us declare it, even if the options are correct
	if options.Name == "" {
		if options.Durable != true || options.Type != Direct {
			return nil, fmt.Errorf("invalid options for default exchange")
		}

		return &Exchange{
			b:    b,
			name: "",
		}, nil
	}
	err := b.mainChan.ExchangeDeclare(
		options.Name,         // name
		string(options.Type), // kind
		options.Durable,      // durable
		!options.Durable,     // autoDelete
		false,                // internal
		false,                // noWait
		nil,                  // args
	)
	if err != nil {
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &Exchange{
		b:    b,
		name: options.Name,
	}, nil
}

func (b *Broker) Queue(options QueueOptions) (*Queue, error) {
	serverdata, err := b.mainChan.QueueDeclare(
		options.Name,      // name
		options.Durable,   // durable
		!options.Durable,  // autoDelete
		options.Exclusive, // exclusive
		false,             // noWait
		nil,               // args
	)
	if err != nil {
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	return &Queue{
		b:                b,
		name:             serverdata.Name,
		consumeExclusive: options.ConsumeExclusive,
		consumerName:     options.ConsumerName,
	}, nil
}

func (b *Broker) Close() error {
	err := b.connection.Close()
	if err != nil {
		return fmt.Errorf("close broker connection: %w", err)
	}

	return nil
}
