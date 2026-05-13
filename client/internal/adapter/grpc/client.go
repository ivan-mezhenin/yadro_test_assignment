package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "dns-manager/proto/dns"
)

type Client struct {
	conn    *grpc.ClientConn
	client  pb.DnsServiceClient
	timeout time.Duration
}

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		client:  pb.NewDnsServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Add(address string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.client.AddDnsServer(ctx, &pb.AddDnsServerRequest{Address: address})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return false, st.Message()
		}
		return false, err.Error()
	}
	return true, ""
}

func (c *Client) Remove(address string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	_, err := c.client.RemoveDnsServer(ctx, &pb.RemoveDnsServerRequest{Address: address})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return false, st.Message()
		}
		return false, err.Error()
	}
	return true, ""
}

func (c *Client) List() (bool, string, []string) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, err := c.client.ListDnsServers(ctx, &pb.ListDnsServersRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return false, st.Message(), nil
		}
		return false, err.Error(), nil
	}
	return true, "", resp.Servers
}
