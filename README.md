# bcamqp

Library for using RabbitMQ/AMQP with Go

## Motivation

This library is a wrapper around the most popular client library: github.com/streadway/amqp

Its interface is designed around the protocol and therefore isn't the most intuitive to use. Additionally, it only serves primitives for you to stick together, without handling the hard bits like reconnecting.

I think a good library should provide enough abstraction and features around a protocol. This is why I built bcamqp.

## Features

- Named parameter objects instead of positional arguments
- Abstraction around protocol primitives
- First crude attempt at reconnection handling
