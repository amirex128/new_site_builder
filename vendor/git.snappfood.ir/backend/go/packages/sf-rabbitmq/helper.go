package sfrabbitmq

import (
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DefaultDial returns a dial function that uses amqp.DefaultDial with the provided timeout
func DefaultDial(timeout time.Duration) func(network, addr string) (net.Conn, error) {
	return amqp.DefaultDial(timeout)
}

// GetLogger returns the global registry's logger or nil if no logger has been set
func GetLogger() Logger {
	return globalRegistry.logger
}
