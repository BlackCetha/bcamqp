package bcamqp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type Exchange struct {
	b    *Broker
	name string
}

func (e *Exchange) Bind(q *Queue, key string) error {
	err := e.b.mainChan.QueueBind(
		q.name, // name
		key,    // key
		e.name, // exchange
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("declare binding: %w", err)
	}

	return nil
}
func (e *Exchange) Unbind(q *Queue, key string) error {
	err := e.b.mainChan.QueueUnbind(
		q.name, // name
		key,    // key
		e.name, // exchange
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("declare queue unbind: %w", err)
	}

	return nil
}
func (e *Exchange) Send(msg Message) error {
	if e.b.autoTimestamp && msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	expiration := ""
	if msg.Expiration > 0 {
		expiration = strconv.FormatInt(msg.Expiration.Milliseconds(), 10)
	}

	exchange := e.name
	if msg.Exchange != "" {
		exchange = msg.Exchange
	}

	dmode := amqp.Persistent
	if msg.Transient {
		dmode = amqp.Transient
	}

	return e.b.mainChan.Publish(
		exchange,
		msg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   msg.ContentType,
			Headers:       msg.Headers,
			Body:          msg.Body,
			Timestamp:     msg.Timestamp,
			CorrelationId: msg.CorrelationID,
			ReplyTo:       msg.ReplyTo,
			Expiration:    expiration,
			DeliveryMode:  dmode,
		},
	)
}
