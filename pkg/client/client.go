package client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Addr string
	Conn *grpc.ClientConn
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &Client{addr, conn}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
