package network

import (
	"net"

	"collectd.org/api"
)

// DefaultService is the default port used by collectd's network plugin.
const DefaultService = "25826"

// Conn is a client connection to a collectd server.
type Conn struct {
	udp    net.Conn
	buffer *Buffer
}

// Dial connects to the collectd server at address. "address" must be a network
// address accepted by net.Dial().
func Dial(address string) (*Conn, error) {
	c, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}

	return &Conn{
		udp:    c,
		buffer: NewBuffer(c),
	}, nil
}

// WriteValueList adds a ValueList to the internal buffer. Data is only written
// to the network when the buffer is full.
func (c *Conn) WriteValueList(vl api.ValueList) error {
	return c.buffer.WriteValueList(vl)
}

// Close closes a connection. You must not use "c" after this call.
func (c *Conn) Close() error {
	if err := c.buffer.Flush(); err != nil {
		return err
	}

	if err := c.udp.Close(); err != nil {
		return err
	}

	c.buffer = nil
	return nil
}