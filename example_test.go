package bcamqp

import "log"

func Example() {
	broker := New(BrokerOptions{
		Encrypted:     true,
		Address:       "localhost",
		User:          "guest",
		Password:      "guest",
		AutoTimestamp: true,
	})

	err := broker.Connect()
	if err != nil {
		log.Fatalf("connect to broker: %v", err)
	}
	defer broker.Close()

	err = broker.DeclareExchange(ExchangeOptions{
		Name:    "example-exchange",
		Durable: false,
		Type:    Direct,
	})
	if err != nil {
		log.Fatalf("declare exchange: %v", err)
	}

	err = broker.DeclareQueue(QueueOptions{
		Name:      "example-queue",
		Durable:   false,
		Exclusive: false,
	})
	if err != nil {
		log.Fatalf("declare queue: %v", err)
	}

	err = broker.DeclareBinding(BindingOptions{
		Exchange:   "example-exchange",
		Queue:      "example-queue",
		RoutingKey: "#",
	})
	if err != nil {
		log.Fatalf("declare binding: %v", err)
	}

	err = broker.Publish(Message{
		Exchange:    "example-exchange",
		Body:        []byte("test"),
		ContentType: "text/plain",
	})
	if err != nil {
		log.Fatalf("publish: %v", err)
	}

	cons, err := broker.Consume(ConsumerOptions{
		Name:      "bcamqp-example",
		Queue:     "example-queue",
		AutoAck:   false,
		Exclusive: false,
	})
	if err != nil {
		log.Fatalf("consume: %v", err)
	}
	defer cons.Close()

	for msg := range cons.Messages() {
		log.Printf(`%v message with routing key "%s": %s`, msg.Timestamp, msg.RoutingKey, msg.Body)

		err = msg.Ack()
		if err != nil {
			log.Printf("ack message: %v", err)
		}
	}
}
