package bcamqp

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// BrokerOptions holds options for broker setup
type BrokerOptions struct {
	Encrypted     bool
	Address       string
	User          string
	Password      string
	AutoTimestamp bool
}

type connectionObserver interface {
	handleConnect()
	handleDisconnect()
}

// Broker wraps the complexities of the AMQP protocol
type Broker struct {
	url                 string
	isReady             bool
	conn                *amqp.Connection
	mainChannelUsed     sync.Mutex
	mainChannel         *amqp.Channel
	autoTimestamp       bool
	connectionObservers map[connectionObserver]struct{}
	reconnectAttempts   int
}

// New returns a new broker instance
func New(options BrokerOptions) *Broker {
	amqpURL := url.URL{
		Scheme: "amqp",
		Host:   options.Address,
		User:   url.UserPassword(options.User, options.Password),
	}
	if options.Encrypted {
		amqpURL.Scheme = "amqps"
	}

	b := Broker{
		url:                 amqpURL.String(),
		autoTimestamp:       options.AutoTimestamp,
		connectionObservers: make(map[connectionObserver]struct{}),
	}
	b.mainChannelUsed.Lock()

	return &b
}

// Connect establishes a connection to the broker
func (b *Broker) Connect() error {
	var err error

	b.conn, err = amqp.Dial(b.url)
	if err != nil {
		return fmt.Errorf("could not connect to broker: %v", err)
	}

	b.mainChannel, err = b.conn.Channel()
	if err != nil {
		return fmt.Errorf("could not open channel to broker: %v", err)
	}

	disconnectChan := make(chan *amqp.Error)
	b.conn.NotifyClose(disconnectChan)
	go b.disconnectHandler(disconnectChan)

	b.mainChannelUsed.Unlock()

	b.isReady = true
	for obs := range b.connectionObservers {
		obs.handleConnect()
	}

	return nil
}

func (b *Broker) disconnectHandler(disconnectChan <-chan *amqp.Error) {
	_, abnormal := <-disconnectChan

	b.mainChannelUsed.Lock()
	b.isReady = false
	for obs := range b.connectionObservers {
		obs.handleDisconnect()
	}

	// Do not reconnect on normal shutdown
	if !abnormal {
		return
	}

	var err error
	var t *time.Timer
	t = time.AfterFunc(time.Duration(b.reconnectAttempts*5)*time.Second, func() {
		err = b.Connect()
		if err == nil {
			return
		}

		t.Reset(time.Duration(b.reconnectAttempts*5) * time.Second)
	})
}

// Message is an AMQP message entity
type Message struct {
	Exchange      string
	RoutingKey    string
	Body          []byte
	Headers       map[string]interface{}
	Timestamp     time.Time
	ContentType   string
	CorrelationID string
	ReplyTo       string
	Expiration    time.Duration
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

// Publish sends a message to the broker
//
// Publishing is blocking concurrency safe
func (b *Broker) Publish(msg Message) error {
	if b.autoTimestamp && msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	expiration := ""
	if msg.Expiration > 0 {
		expiration = strconv.FormatInt(msg.Expiration.Milliseconds(), 10)
	}

	b.mainChannelUsed.Lock()
	defer b.mainChannelUsed.Unlock()

	return b.mainChannel.Publish(
		msg.Exchange,
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
		},
	)
}

// Close gracefully shuts down the broker connection
func (b *Broker) Close() {
	b.isReady = false

	b.conn.Close()
}

func (b *Broker) subscribe(obs connectionObserver) {
	b.connectionObservers[obs] = struct{}{}
}

func (b *Broker) unsubscribe(obs connectionObserver) {
	delete(b.connectionObservers, obs)
}
