package clients

import (
	"TGNotification/pkg/api"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DjangoClient struct {
	client api.NotificationSystemClient
	conn   *grpc.ClientConn
}

func NewDjangoClient(addr string) (*DjangoClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("clients.NewDjangoClient, failed to connect to django grpc %v", err)
	}
	client := api.NewNotificationSystemClient(conn)

	return &DjangoClient{
		client: client,
		conn:   conn}, nil
}

func (c *DjangoClient) Close() error {
	return c.conn.Close()
}

func (c *DjangoClient) LinkUser(ctx context.Context, token string, chatId int64, uname string) (bool, string, error) {
	req := &api.LinkRequest{
		Token:    token,
		ChatId:   chatId,
		Username: uname,
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	resp, err := c.client.LinkUserTelegram(ctxTimeout, req)
	if err != nil {
		return false, "", fmt.Errorf("clients.DjangoClient.LinkUser, grpc call failed %v", err)
	}

	return resp.Success, resp.Message, nil
}
