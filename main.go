package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/coreos/dex/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	dexCAPath = flag.String("ca-cert", "/etc/dex/grpc.crt", "Path to CA certificate")
	dexHost   = flag.String("dex-host", "127.0.0.1:5557", "The host:port to dex's gRPC port")
)

func newDexClient(hostAndPort, caPath string) (api.DexClient, error) {
	creds, err := credentials.NewClientTLSFromFile(caPath, "")
	if err != nil {
		return nil, fmt.Errorf("load dex cert: %v", err)
	}

	conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return api.NewDexClient(conn), nil
}

func main() {
	flag.Parse()

	client, err := newDexClient(*dexHost, *dexCAPath)
	if err != nil {
		log.Fatalf("failed creating dex client: %v ", err)
	}

	req := &api.VersionReq{}
	rsp, err := client.GetVersion(context.Background(), req)
	if err != nil {
		log.Fatalf("failed querying for version: %v", err)
	}

	log.Printf("got version: %+v", rsp)

	// req := &api.CreateClientReq{
	// 	Client: &api.Client{
	// 		Id:           "example-app",
	// 		Name:         "Example App",
	// 		Secret:       "ZXhhbXBsZS1hcHAtc2VjcmV0",
	// 		RedirectUris: []string{"http://127.0.0.1:5555/callback"},
	// 	},
	// }

	// if _, err := client.CreateClient(context.TODO(), req); err != nil {
	// 	log.Fatalf("failed creating oauth2 client: %v", err)
	// }
}
