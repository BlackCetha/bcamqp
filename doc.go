// Package bcamqp is a RabbitMQ/AMQP client library.
//
// AMQP is a messaging protocol, RabbitMQ is a server implementation of that protocol.
// Read more about them here: https://www.rabbitmq.com/tutorials/amqp-concepts.html
//
// bcamqp allows your program to send and receive messages via the AMQP protocol with relative ease.
// It tries to abstract away the complexities and provide a clear interface following Go's usual style.
//
// There is also basic reconnection logic in place. While the connection is down, the libary will block.
package bcamqp
