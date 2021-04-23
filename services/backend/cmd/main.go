package main

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/DaniilOr/spamer/services/backend/cmd/app"
	"github.com/DaniilOr/spamer/services/backend/pkg/auth"
	"github.com/DaniilOr/spamer/services/backend/pkg/classifier"
	"github.com/DaniilOr/spamer/services/backend/pkg/spam"
	"github.com/go-chi/chi"
	"go.opencensus.io/trace"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultPort               = "9999"
	defaultHost               = "0.0.0.0"
	defaultAuthURL            = "auth:8080"
	defaultClassifierAPIURL = "classifier:9090"
	defaultSpamAPIURL = "spam:8888"
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

	authURL, ok := os.LookupEnv("APP_AUTH_URL")
	if !ok {
		authURL = defaultAuthURL
	}

	classifierURL, ok := os.LookupEnv("APP_CLASSIFIER_URL")
	if !ok {
		classifierURL = defaultClassifierAPIURL
	}
	spamURL, ok := os.LookupEnv("APP_SPAM_URL")
	if !ok {
		spamURL = defaultSpamAPIURL
	}
	err := InitJaeger("backed")
	if err != nil{
		log.Println(err)
		os.Exit(1)
	}
	if err := execute(net.JoinHostPort(host, port), authURL, classifierURL, spamURL); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, authURL string, classifierURL string, spamURL string) error {
	authSvc, err := auth.Init(authURL)
	if err != nil{
		return err
	}

	mux := chi.NewRouter()
	classifierSVC, err := classifier.Init(classifierURL)
	if err != nil{
		return err
	}
	spamSVC, err := spam.Init(spamURL)
	if err != nil{
		return err
	}
	application := app.NewServer(authSvc, classifierSVC, spamSVC, mux)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
