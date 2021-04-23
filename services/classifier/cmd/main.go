package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/DaniilOr/spamer/services/classifier/cmd/app"
	serverPb "github.com/DaniilOr/spamer/services/classifier/pkg/server"
	"github.com/DaniilOr/spamer/services/classifier/pkg/URLC"
	"github.com/DaniilOr/spamer/services/classifier/pkg/SMSC"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const (
	defaultPort = "8080"
	defaultHost = "0.0.0.0"
	defaultURL  = ""
)
func InitJaeger(serviceName string) error{
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "jaeger:6831",
		Process: jaeger.Process{
			ServiceName: serviceName,
			Tags: []jaeger.Tag{
				jaeger.StringTag("hostname", "localhost"),
			},
		},
	})
	if err != nil {
		return err
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
	return nil
}
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
	err := InitJaeger("auth")
	if err != nil{
		log.Println(err)
		os.Exit(1)
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
