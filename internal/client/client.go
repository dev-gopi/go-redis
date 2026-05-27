package client

import (
	"bufio"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	ID            string
	Conn          net.Conn
	Reader        *bufio.Reader
	Authenticated bool
	Username      string
	SelectedDB    int
}

func NewClient(conn net.Conn) *Client {

	return &Client{
		ID:            uuid.NewString(),
		Conn:          conn,
		Reader:        bufio.NewReader(conn),
		Authenticated: false,
		Username:      "",
		SelectedDB:    0,
	}
}
