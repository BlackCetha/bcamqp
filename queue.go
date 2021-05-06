package bcamqp

import (
	"fmt"
	"time"

	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/streadway/amqp"
)

type Queue struct {
	b                *Broker
	name             string
	err              error
	channel          adaptor.Channel
	msg              *Message
	msgsChan         <-chan amqp.Delivery
	consumeExclusive bool
	consumerName     string
}

func (q *Queue) Next() bool {
	q.err = nil

	if q.channel == nil {
		if !q.beginConsuming() {
			return false
		}
	}

	m, ok := <-q.msgsChan
	if !ok {
		return false
	}

	q.msg = &Message{
		Exchange:      m.Exchange,
		RoutingKey:    m.RoutingKey,
		Body:          m.Body,
		ContentType:   m.ContentType,
		CorrelationID: m.CorrelationId,
		Headers:       m.Headers,
		ReplyTo:       m.ReplyTo,
		Timestamp:     m.Timestamp,
		ackFunc:       m.Ack,
		rejectFunc:    m.Reject,
		nackFunc:      m.Nack,
	}

	return true
}

func (q *Queue) NextWithTimeout(timeout time.Duration) bool {
	q.err = nil

	if q.channel == nil {
		if !q.beginConsuming() {
			return false
		}
	}

	timeoutChan := time.After(timeout)

	var m amqp.Delivery
	var ok bool
	select {
	case m, ok = <-q.msgsChan:
	case <-timeoutChan:
		return false
	}

	if !ok {
		return false
	}

	q.msg = &Message{
		Exchange:      m.Exchange,
		RoutingKey:    m.RoutingKey,
		Body:          m.Body,
		ContentType:   m.ContentType,
		CorrelationID: m.CorrelationId,
		Headers:       m.Headers,
		ReplyTo:       m.ReplyTo,
		Timestamp:     m.Timestamp,
		ackFunc:       m.Ack,
		rejectFunc:    m.Reject,
		nackFunc:      m.Nack,
	}

	return true
}

func (q *Queue) Message() *Message {
	return q.msg
}

func (q *Queue) Error() error {
	return q.err
}

func (q *Queue) Delete() error {
	_, err := q.b.mainChan.QueueDelete(
		q.name, // name
		false,  // ifUnused
		false,  // ifEmpty
		false,  // noWait
	)
	if err != nil {
		return fmt.Errorf("delete queue: %w", err)
	}

	return nil
}

func (q *Queue) Close() error {
	if q.channel == nil {
		return nil
	}

	err := q.channel.Cancel(
		q.consumerName, // consumer
		false,          // noWait
	)
	if err != nil {
		return fmt.Errorf("stop consumer: %w", err)
	}

	err = q.channel.Close()
	if err != nil {
		return fmt.Errorf("close channel: %w", err)
	}

	return nil
}

func (q *Queue) beginConsuming() bool {
	var err error
	q.channel, err = q.b.connection.Channel()
	if err != nil {
		q.err = fmt.Errorf("open channel: %w", err)
		return false
	}

	q.msgsChan, err = q.channel.Consume(
		q.name,             // queue
		q.consumerName,     // consumer
		false,              // autoAck
		q.consumeExclusive, // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // args
	)
	if err != nil {
		// an error here is most likely a channel-level exception and will
		// close the channel anyway, make sure its closed and carry on
		_ = q.channel.Close()
		q.channel = nil

		q.err = fmt.Errorf("begin consuming: %w", err)
		return false
	}

	return true
}
