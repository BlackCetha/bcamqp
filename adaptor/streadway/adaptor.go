// Package streadway is an adaptor for the de-facto standard AMQP library.
package streadway

import (
	"github.com/blackcetha/bcamqp/adaptor"
	"github.com/streadway/amqp"
)

// make mistakes compile time errors
var _ adaptor.Server = Server{}
var _ adaptor.Connection = &Connection{}
var _ adaptor.Channel = &amqp.Channel{}

type Server struct{}

func (s Server) Dial(url string) (adaptor.Connection, error) {
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
