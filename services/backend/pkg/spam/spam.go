package spam

import (
	"context"
	serverPb "github.com/DaniilOr/spamer/services/spam/pkg/server"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

type Service struct{
	client serverPb.SpamerClient

}

func Init(addr string) (*Service, error){
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return  nil, err
	}
	client :=  serverPb.NewSpamerClient(conn)
	server := Service{client: client}
	return &server, nil
}

func (s*Service) Spam(ctx context.Context, url string) (string, error) {
	ctx, span := trace.StartSpan(ctx, "route: token")
	defer span.End()
	response, err := s.client.Spam(ctx, &serverPb.Target{Url: url})
	if err != nil{
		return "", err
	}
	return response.Verdict, nil
}