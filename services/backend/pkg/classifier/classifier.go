package classifier

import (
	"context"
	serverPb "github.com/DaniilOr/spamer/services/classifier/pkg/server"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"log"
)

type Service struct{
	client serverPb.ClassifierClient

}

func Init(addr string) (*Service, error){
	log.Println("Enter Init classifier")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	log.Println("Init classifier")
	if err != nil {
		return  nil, err
	}
	client :=  serverPb.NewClassifierClient(conn)
	server := Service{client: client}
	return &server, nil
}

func (s*Service) CheckURL(ctx context.Context, url string) (string, error) {
	log.Println("Inside classifier")
	ctx, span := trace.StartSpan(ctx, "route: url")
	defer span.End()
	log.Println("CHeck URL inside backend")
	response, err := s.client.CheckURL(ctx, &serverPb.URLReq{Url: url})
	if err != nil{
		return "", err
	}
	return response.Verdict, nil
}
func (s*Service) CheckSMS(ctx context.Context, sms string) (string, error) {
	ctx, span := trace.StartSpan(ctx, "route: sms")
	defer span.End()
	response, err := s.client.CheckSMS(ctx, &serverPb.SMSReq{Sms: sms})
	if err != nil{
		return "", err
	}
	return response.Verdict, nil
}