package bcamqp

// ExchangeType specifies the type of the exchange to be created
type ExchangeType string

// Lists the different exchange types
const (
	Direct  ExchangeType = "direct"
	Fanout  ExchangeType = "fanout"
	Topic   ExchangeType = "topic"
	Headers ExchangeType = "headers"
)

// ExchangeOptions holds options for exchange creation
type ExchangeOptions struct {
	Name    string
	Type    ExchangeType
	Durable bool
}

// DeclareExchange makes sure that there is an exchange with the specified
// properties on the broker
//
// It returns an error when a exchange with the specified name already exists,
// but differs
func (b *Broker) DeclareExchange(options ExchangeOptions) error {
	return b.mainChannel.ExchangeDeclare(
		options.Name,         // name
		string(options.Type), // kind
		options.Durable,      // durable
		!options.Durable,     // autoDelete
		false,                // internal
		false,                // noWait
		nil,
	)
}

// QueueOptions holds options for queue creation
type QueueOptions struct {
	Name      string
	Durable   bool
	Exclusive bool
}

// DeclareQueue makes sure that there is a queue with the specified properties
// on the broker
//
// It returns an error when a queue with the specified name already exisits,
// but differs
func (b *Broker) DeclareQueue(options QueueOptions) error {
	_, err := b.mainChannel.QueueDeclare(
		options.Name,      // name
		options.Durable,   // durable
		!options.Durable,  // autoDelete
		options.Exclusive, // exclusive
		false,             // noWait
		nil,
	)

	return err
}

// BindingOptions holds options for binding creation
type BindingOptions struct {
	Name       string
	RoutingKey string
	Exchange   string
}

// DeclareBinding creates a binding between an exchange and a queue
func (b *Broker) DeclareBinding(options BindingOptions) error {
	return b.mainChannel.QueueBind(
		options.Name,       // name
		options.RoutingKey, // key
		options.Exchange,   // exchange
		false,              // noWait
		nil,
	)
}
