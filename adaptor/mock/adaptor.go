package mock

import (
	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/streadway/amqp"
)

var _ adaptor.Server = Server{}
var _ adaptor.Connection = &Connection{}
var _ adaptor.Channel = &Channel{}

type Server struct{}

func (s Server) Dial(url string) (adaptor.Connection, error) {
	return &Connection{}, nil
}

type Connection struct{}

func (c *Connection) Channel() (adaptor.Channel, error) {
	return &Channel{}, nil
}

func (c *Connection) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return receiver
}

func (c *Connection) Close() error {
	return nil
}

type Channel struct{}

func (c *Channel) Cancel(consumer string, noWait bool) error {
	return nil
}

func (c *Channel) Close() error {
	return nil
}

func (c *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (c *Channel) ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error {
	return nil
}

func (c *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return nil
}

func (c *Channel) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return nil
}

func (c *Channel) ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error {
	return nil
}

func (c *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return nil
}

func (c *Channel) Qos(prefetchCount, prefetchSize int, global bool) error {
	return nil
}

func (c *Channel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return nil
}

func (c *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (c *Channel) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	return 0, nil
}

func (c *Channel) QueueUnbind(name, key, exchange string, args amqp.Table) error {
	return nil
}
