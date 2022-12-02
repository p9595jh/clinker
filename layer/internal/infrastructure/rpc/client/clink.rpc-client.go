package rpcclient

import (
	"context"
	"layer/internal/infrastructure/rpc/proto/clink"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClinkRpcClient struct {
	client clink.ClinkClient
}

func NewClinkRpcClient(rpcUrl string) *ClinkRpcClient {
	conn, err := grpc.Dial(rpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return &ClinkRpcClient{
		client: clink.NewClinkClient(conn),
	}
}

func (c *ClinkRpcClient) Confirm(kind clink.Kind, id string) (*clink.ConfirmResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return c.client.Confirm(ctx, &clink.ConfirmRequest{
		Kind: kind,
		Id:   id,
	})
}

func (c *ClinkRpcClient) ConfirmError(kind clink.Kind, err error) (*clink.ConfirmResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return c.client.Confirm(ctx, &clink.ConfirmRequest{
		Kind:  kind,
		Error: err.Error(),
	})
}
