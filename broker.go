package bcamqp

import (
	"fmt"
	"strings"

	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/streadway/amqp"
)

type Broker struct {
	connection          adaptor.Connection
	mainChan            adaptor.Channel
	disconnectListeners []disconnectListener
	autoTimestamp       bool
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

func (b *Broker) unregisterDisconnectListener(listener disconnectListener) {
	for i := 0; i < len(b.disconnectListeners); i++ {
		if b.disconnectListeners[i] == listener {
			b.disconnectListeners = append(b.disconnectListeners[:i], b.disconnectListeners[i+1:]...)
		}
	}
}

func (b *Broker) GetExchange(options ExchangeOptions) (*Exchange, error) {
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

func (b *Broker) GetQueue(options QueueOptions) (*Queue, error) {
	serverdata, err := b.mainChan.QueueDeclare(
		options.Name,      //
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
	var err error
	var errs []error
	for _, l := range b.disconnectListeners {
		err = l.onDisconnect()
		if err != nil {
			errs = append(errs, err)
		}
	}

	err = b.mainChan.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("close main channel: %w", err))
	}

	err = b.connection.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("close broker connection: %w", err))
	}

	if len(errs) == 0 {
		return nil
	}

	var out strings.Builder
	for i := 0; i < len(errs)-1; i++ {
		out.WriteString(errs[i].Error())
		out.WriteString(", ")
	}
	out.WriteString(errs[len(errs)-1].Error())

	return fmt.Errorf(out.String())
}
