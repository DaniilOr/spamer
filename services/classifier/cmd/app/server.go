package app

import (
	"context"
	"github.com/DaniilOr/spamer/services/classifier/pkg/SMSC"
	"github.com/DaniilOr/spamer/services/classifier/pkg/URLC"
	serverPb "github.com/DaniilOr/spamer/services/classifier/pkg/server"
	"go.opencensus.io/trace"

	"log"
)

type Server struct {
	Sms *SMSC.Service
	Url *URLC.Service
	ctx context.Context
}

func NewServer(sms *SMSC.Service, url *URLC.Service, ctx context.Context) *Server {
	return &Server{Sms: sms, Url: url, ctx: ctx}
}

func (s *Server) CheckURL(ctx context.Context, request *serverPb.URLReq) ( * serverPb.URLResp, error) {
	ctx, span := trace.StartSpan(ctx, "route: token")
	defer span.End()
	res, err := s.Url.CheckURL(request.Url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.URLResp{Verdict: res}
	return &response, nil
}

func (s *Server) CheckSMS(ctx context.Context, request *serverPb.SMSReq) ( * serverPb.SMSResp, error) {
	ctx, span := trace.StartSpan(ctx, "route: token")
	defer span.End()
	res, err := s.Sms.CheckSMS(request.Sms)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.SMSResp{Verdict: res}
	return &response, nil
}
