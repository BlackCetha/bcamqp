package streadway

import (
	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/streadway/amqp"
)

var _ adaptor.Server = Streadway{}
var _ adaptor.Connection = &Connection{}
var _ adaptor.Channel = &amqp.Channel{}

type Streadway struct{}

func (s Streadway) Dial(url string) (adaptor.Connection, error) {
	conn, err := amqp.Dial(url)
	return &Connection{conn: conn}, err
}

type Connection struct {
	conn *amqp.Connection
}

func (c *Connection) Channel() (adaptor.Channel, error) {
	return c.conn.Channel()
}

func (c *Connection) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return c.conn.NotifyClose(receiver)
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
