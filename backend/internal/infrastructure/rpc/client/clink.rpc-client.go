package rpcclient

import (
	"clinker-backend/internal/infrastructure/database/entity"
	"clinker-backend/internal/infrastructure/rpc/proto/clink"
	"context"

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

func (c *ClinkRpcClient) CreateVestige(vestige *entity.Vestige) (*clink.TxHashResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return c.client.VestigeCreate(ctx, &clink.VestigeCreateRequest{
		Parent:  vestige.Parent,
		Head:    vestige.Head,
		Title:   vestige.Title,
		Content: vestige.Content,
		Hit:     vestige.Hit,
		UserId:  vestige.UserId,
		Address: vestige.User.Address,
	})
}

func (c *ClinkRpcClient) CreateAppraisal(appraisal *entity.Appraisal) (*clink.TxHashResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return c.client.AppraisalCreate(ctx, &clink.AppraisalCreateRequest{
		Value:     appraisal.Value,
		VestigeId: appraisal.VestigeId,
		NextId:    appraisal.NextId,
		UserId:    appraisal.UserId,
		Address:   appraisal.User.Address,
	})
}
