package discovery

import (
	"context"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func TestGetService(t *testing.T) {
	ctx := context.Background()
	client, err := clientv3.New(clientv3.Config{
		Context:   ctx,
		Endpoints: []string{"127.0.0.1:2379"},
		// DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	service := NewGoMicro(client)

	r, err := service.GetService(ctx, "servicename")
	if err != nil {
		t.Fatal(err)
	}

	if len(r) == 0 {
		t.Errorf("service not empty")
	}
}
