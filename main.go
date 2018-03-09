package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/dex/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func newDexClient(hostAndPort, caPath string) (api.DexClient, error) {
	creds, err := credentials.NewClientTLSFromFile(caPath, "")
	if err != nil {
		return nil, fmt.Errorf("load dex cert: %v", err)
	}

	conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dail: %v", err)
	}
	return api.NewDexClient(conn), nil
}

func main() {
	client, err := newDexClient("127.0.0.1:5557", "/etc/dex/grpc.crt")
	if err != nil {
		log.Fatalf("failed creating dex client: %v ", err)
	}

	req := &api.CreateClientReq{
		Client: &api.Client{
			Id:           "example-app",
			Name:         "Example App",
			Secret:       "ZXhhbXBsZS1hcHAtc2VjcmV0",
			RedirectUris: []string{"http://127.0.0.1:5555/callback"},
		},
	}

	if _, err := client.CreateClient(context.TODO(), req); err != nil {
		log.Fatalf("failed creating oauth2 client: %v", err)
	}
}
