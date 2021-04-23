package app

import (
	"context"
	serverPb "github.com/DaniilOr/spamer/services/spam/pkg/server"
	"go.opencensus.io/trace"
	"log"
	"spamer/services/spam/pkg/spam"
)

type Server struct {
	spam *spam.Service
	ctx context.Context
}

func NewServer(spam *spam.Service, ctx context.Context) *Server {
	return &Server{spam: spam, ctx: ctx }
}

func (s *Server) Spam(ctx context.Context, request *serverPb.Target) ( * serverPb.Response, error) {
	ctx, span := trace.StartSpan(ctx, "route: token")
	defer span.End()
	res, err := s.spam.Spam(ctx, request.url, request.interval, request.numStreams)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.Response{verdict: res}
	return &response, nil
}