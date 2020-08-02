package rediscli

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Client struct {
	conn  *net.Conn
	Debug bool
}

func New(address string) (*Client, error) {
	var err error
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	client := &Client{
		&conn,
		false,
	}
	return client, nil
}

func (c *Client) Auth(password string) error {
	_, err := c.send(fmt.Sprintf("AUTH %s", password))
	return err
}

func (c *Client) Get(key string) (interface{}, error) {
	return c.send(fmt.Sprintf("GET %s", key))
}

func (c *Client) Ping() (interface{}, error) {
	return c.send("PING")
}

func (c *Client) send(cmd string) (interface{}, error) {
	if c.Debug {
		fmt.Println("<-", cmd)
	}

	fmt.Fprintf(*c.conn, fmt.Sprintf("%s\r\n", cmd))

	data, err := c.readLine()

	if err != nil {
		return nil, err
	}

	switch data[0] {
	case '-':
		// Error
		return nil, errors.New(string(data[1 : len(data)-2]))
	case '+':
		// Simple string
		return string(data[1 : len(data)-2]), nil
	case '$':
		// Empty result
		if data[1] == '-' {
			return nil, nil
		}

		// Blob string
		size := binary.BigEndian.Uint64(data[1 : len(data)-2])
		results, err := c.readBytes(size)
		if err != nil {
			return nil, err
		}
		return string(results), nil
	default:
		return string(data), nil
	}
}

func (c *Client) readLine() ([]byte, error) {
	reader := bufio.NewReader(*c.conn)
	buf, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	if c.Debug {
		fmt.Println("->", string(buf))
	}

	size := len(buf)
	if size <= 2 || buf[size-2] != '\r' || buf[size-1] != '\n' {
		return nil, fmt.Errorf("invalid reply: %q", buf)
	}
	return buf, nil
}

func (c *Client) readBytes(size uint64) ([]byte, error) {
	buf := make([]byte, size+1)
	reader := bufio.NewReader(*c.conn)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}

	if c.Debug {
		fmt.Println("->", string(buf))
	}
	if buf[size] != '\r' || buf[size+1] != '\n' {
		return nil, fmt.Errorf("invalid reply: %q", buf)
	}
	return buf, nil
}
