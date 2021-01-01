package bcamqp

import (
	"sync"

	"github.com/streadway/amqp"
)

// Consumer gets messages from the broker
type Consumer struct {
	amqpChanUsed sync.Mutex
	amqpChan     *amqp.Channel
	messages     chan Message
	options      ConsumerOptions
	broker       *Broker
}

// Messages returns a channel from which incoming messages can be read
func (c *Consumer) Messages() <-chan Message {
	return c.messages
}

// Close gracefully shuts down the consumer
//
// Trying to read from a closed consumer
func (c *Consumer) Close() error {
	// Our protocol library (github.com/streadway/amqp) will
	// synchronize this call itself. We need to call this
	// prior to acquiring the mutex via handleDisconnect()
	// because on clean shutdown the message loop from
	// startConsuming will still be running and holding
	// the mutex. This will run the last batch of messages
	// through it and then close the delivery channel itself.
	// This in turn will allow startConsuming to return and
	// release the mutex.
	c.amqpChan.Cancel(c.options.Name, false)
	
	// Will acquire the mutex
	c.handleDisconnect()
	
	return c.amqpChan.Close()
}

func (c *Consumer) startConsuming() {
	c.amqpChanUsed.Lock()
	defer c.amqpChanUsed.Unlock()

	deliveries, err := c.amqpChan.Consume(
		c.options.Queue,
		c.options.Name,
		c.options.AutoAck,
		c.options.Exclusive,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	for msg := range deliveries {
		c.messages <- Message{
			Exchange:      msg.Exchange,
			RoutingKey:    msg.RoutingKey,
			Body:          msg.Body,
			ContentType:   msg.ContentType,
			CorrelationID: msg.CorrelationId,
			Headers:       msg.Headers,
			ReplyTo:       msg.ReplyTo,
			Timestamp:     msg.Timestamp,
			ackFunc:       msg.Ack,
			rejectFunc:    msg.Reject,
			nackFunc:      msg.Nack,
		}
	}
}

func (c *Consumer) handleConnect() {
	var err error
	c.amqpChan, err = c.broker.conn.Channel()
	if err != nil {
		return
	}

	c.amqpChanUsed.Unlock()

	go c.startConsuming()
}

func (c *Consumer) handleDisconnect() {
	c.amqpChanUsed.Lock()
}

// ConsumerOptions holds options for consumer setup
type ConsumerOptions struct {
	Name      string // application-defined, e.g. executable name
	Queue     string
	AutoAck   bool
	Exclusive bool
}

// Consume starts a new consumer
func (b *Broker) Consume(options ConsumerOptions) *Consumer {
	c := &Consumer{
		options:  options,
		broker:   b,
		messages: make(chan Message),
	}
	c.amqpChanUsed.Lock()

	if b.isReady {
		c.handleConnect()
	}

	b.subscribe(c)

	return c
}
