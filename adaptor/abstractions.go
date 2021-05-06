// Package adoptor provides abstraction of AMQP messaging libraries.
// It is primarily used to introduce a mock during testing.
package adaptor

import "github.com/streadway/amqp"

// Server wraps the library function used to connect to a broker
// in a struct.
type Server interface {
	Dial(url string) (Connection, error)
}

// Connection wraps a broker connection.
type Connection interface {
	Channel() (Channel, error)
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
	// Close needs to relay to all channels and consumers
	Close() error
}

// Channel wraps the operations of a protocol primitive "channel".
type Channel interface {
	Cancel(consumer string, noWait bool) error
	Close() error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
	QueueUnbind(name, key, exchange string, args amqp.Table) error
}
