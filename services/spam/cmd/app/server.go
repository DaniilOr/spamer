package app

import (
	"context"
	serverPb "github.com/DaniilOr/spamer/services/spam/pkg/server"
	"github.com/DaniilOr/spamer/services/spam/pkg/spam"
	"go.opencensus.io/trace"
	"log"
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
	res, err := s.spam.Spam(request.Url, request.Interval, request.NumStreams)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.Response{Verdict: res}
	return &response, nil
}