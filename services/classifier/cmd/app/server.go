package app

import (
	"context"
	"go.opencensus.io/trace"
	serverPb "github.com/DaniilOr/spamer/services/classifier/pkg/server"
	"github.com/DaniilOr/spamer/services/classifier/pkg/SMSC"
	"github.com/DaniilOr/spamer/services/classifier/pkg/URLC"

	"log"
)

type Server struct {
	Sms *SMSC.Service
	Url *URLC.Service
	ctx context.Context
}

func NewServer(sms *SMSC.Service, url *URLC.Service, ctx context.Context) *Server {
	return &Server{SMSC: sms, URLC: urk, ctx: ctx}
}

func (s *Server) Token(ctx context.Context, request *serverPb.TokenRequest) ( * serverPb.TokenResponse, error) {
	ctx, span := trace.StartSpan(ctx, "route: token")
	defer span.End()
	token, exp, err := s.authSvc.Login(ctx, request.Login, request.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.TokenResponse{Token: token, Expire: exp}
	return &response, nil
}

// Доступно всем
func (s *Server) Id (ctx context.Context, request *serverPb.IdRequest) (*serverPb.IdResponse, error) {
	userID, err := s.authSvc.UserID(ctx, request.Token)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	response := serverPb.IdResponse{UserId: userID}
	return &response, nil
}
