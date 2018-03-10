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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client, err := newDexClient(*dexHost, *dexCAPath, *clientCertPath, *clientKeyPath)
	if err != nil {
		return fmt.Errorf("failed creating dex client: %v ", err)
	}

	if err := getVersion(client); err != nil {
		return err
	}

	return nil
}

func createClient(client api.DexClient) error {
	req := &api.CreateClientReq{
		Client: &api.Client{
			Id:           "example-app",
			Name:         "Example App",
			Secret:       "ZXhhbXBsZS1hcHAtc2VjcmV0",
			RedirectUris: []string{"http://127.0.0.1:5555/callback"},
		},
	}

	if _, err := client.CreateClient(context.TODO(), req); err != nil {
		fmt.Errorf("failed creating oauth2 client: %v", err)
	}

	return nil
}

func getVersion(client api.DexClient) error {
	req := &api.VersionReq{}
	rsp, err := client.GetVersion(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed querying for version: %v", err)
	}

	log.Printf("got version: %+v", rsp)
	return nil
}
