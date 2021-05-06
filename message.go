package bcamqp

import "time"

// Message is an AMQP message entity
//
// These fields are part of the AMQP standard or RabbitMQ extensions
type Message struct {
	Exchange      string
	RoutingKey    string
	Body          []byte
	Headers       map[string]interface{}
	Timestamp     time.Time // application-defined, may be set to any value
	ContentType   string
	CorrelationID string
	ReplyTo       string
	Expiration    time.Duration
	Transient     bool // store message in-mem only
	ackFunc       func(bool) error
	rejectFunc    func(bool) error
	nackFunc      func(bool, bool) error
}

// Ack acknowledges the message, implying that it was processed correctly
// and completely
func (m *Message) Ack() error {
	return m.ackFunc(false)
}

// Reject is used to tell the server that this client is unable to handle
// this message, optionally requeueing it
func (m *Message) Reject(requeue bool) error {
	return m.rejectFunc(requeue)
}

// Nack is used to tell the server that this client is not willing to handle
// this message, optionally requeueing it for other consumers
func (m *Message) Nack(requeue bool) error {
	return m.nackFunc(false, requeue)
}
