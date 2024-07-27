package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) SelectService(service string) error {
	_, err := fmt.Fprintf(c.conn, "%s\n", service)
	if err != nil {
		return err
	}

	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return err
	}

	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "Error") {
		return fmt.Errorf(response)
	}

	return nil
}

func (c *Client) SendCommand(command string) (string, error) {
	_, err := fmt.Fprintf(c.conn, "%s\n", command)
	if err != nil {
		return "", err
	}

	response, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}
