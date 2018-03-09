package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/coreos/dex/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	clientCertPath = flag.String("client-cert", "", "Path to client certificate")
	clientKeyPath  = flag.String("client-key", "", "Path to client key")

	dexCAPath = flag.String("ca-cert", "/etc/dex/grpc.crt", "Path to CA certificate")
	dexHost   = flag.String("dex-host", "127.0.0.1:5557", "The host:port to dex's gRPC port")
)

func newDexClient(hostAndPort, caPath, clientCertPath, clientKeyPath string) (api.DexClient, error) {
	cpool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CA certificate: %v", err)
	}

	if !cpool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	clientCreds, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("invalid client credentials: %v", err)
	}

	clientTLS := &tls.Config{
		RootCAs:      cpool,
		Certificates: []tls.Certificate{clientCreds},
	}

	creds := credentials.NewTLS(clientTLS)
	conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	return api.NewDexClient(conn), nil
}

func main() {
	flag.Parse()

	client, err := newDexClient(*dexHost, *dexCAPath, *clientCertPath, *clientKeyPath)
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
