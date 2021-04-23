package main

import (
	"context"
	"github.com/DaniilOr/spamer/services/classifier/cmd/app"
	"github.com/DaniilOr/spamer/services/classifier/pkg/SMSC"
	"github.com/DaniilOr/spamer/services/classifier/pkg/URLC"
	serverPb "github.com/DaniilOr/spamer/services/classifier/pkg/server"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"net"
	"os"
)

const (
	defaultPort = "9090"
	defaultHost = "0.0.0.0"
	defaultURL  = "mlurl:5000/classify/url"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	url, ok := os.LookupEnv("ML_URL")
	if !ok {
		url = defaultURL
	}
	if err := execute(net.JoinHostPort(host, port), url); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, url string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ctx := context.Background()
	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	Smss := SMSC.NewService(url)
	Urls := URLC.NewService(url)
	server := app.NewServer(Smss, Urls, ctx)
	serverPb.RegisterClassifierServer(grpcServer, server)
	return grpcServer.Serve(listener)
}
