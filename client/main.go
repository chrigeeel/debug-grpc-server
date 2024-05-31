package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	pb "github.com/chrigeeel/debug-grpc-server/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	SERVER_ENDPOINT = "208.91.106.21:29381"
)

func main() {
	proxyURL, _ := url.Parse("http://HJUW814842:FHJKL027@204.52.120.1:6344") //proxy here

	customDialer := func(ctx context.Context, addr string) (net.Conn, error) {
		proxyConn, err := (&http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}).DialContext(ctx, "tcp", addr)
		if err != nil {
			return nil, err
		}

		return proxyConn, nil
	}

	//connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithContextDialer(customDialer),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	conn, err := grpc.DialContext(context.Background(), SERVER_ENDPOINT, opts...)
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewIPServiceClient(conn)

	log.Println(client.GetIP(context.TODO(), nil))
}
