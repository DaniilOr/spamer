package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DaniilOr/spamer/services/classifier/pkg/SMSC"
	"github.com/DaniilOr/spamer/services/classifier/pkg/URLC"
	serverPb "github.com/DaniilOr/spamer/services/classifier/pkg/server"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)
const cacheTimeout = 50 * time.Millisecond
type Server struct {
	Sms *SMSC.Service
	Url *URLC.Service
	ctx context.Context
	cache *redis.Pool
}
type CacheResp struct {
	Url string `json:"url" bson:"url"`
}
func NewServer(sms *SMSC.Service, url *URLC.Service, cache *redis.Pool, ctx context.Context) *Server {
	return &Server{Sms: sms, Url: url, ctx: ctx, cache: cache}
}

func (s *Server) CheckURL(ctx context.Context, request *serverPb.URLReq) ( * serverPb.URLResp, error) {
	log.Println("Enter Check url inside")
	var r CacheResp
	if cached, err := s.FromCache(ctx, fmt.Sprintf("urls:%s", request.Url)); err == nil {
		log.Printf("Got from cache: %s", cached)
		err = json.Unmarshal(cached, r)
		if err != nil{
			return nil, err
		}
		return &serverPb.URLResp{Verdict: r.Url}, nil
	}
	res, err := s.Url.CheckURL(request.Url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	r.Url = res
	data, err := json.Marshal(r)
	if err != nil{
		return nil, err
	}
	err = s.ToCache(ctx, fmt.Sprintf("urls:%s", res), data)
	response := serverPb.URLResp{Verdict: res}
	return &response, nil
}

func (s *Server) CheckSMS(ctx context.Context, request *serverPb.SMSReq) ( * serverPb.SMSResp, error) {
	res, err := s.Sms.CheckSMS(request.Sms)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	response := serverPb.SMSResp{Verdict: res}
	return &response, nil
}

func (s *Server) FromCache(ctx context.Context, key string) ([]byte, error) {
	conn, err := s.cache.GetContext(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	reply, err := redis.DoWithTimeout(conn, cacheTimeout, "GET", key)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	value, err := redis.Bytes(reply, err)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return value, err
}

func (s *Server) ToCache(ctx context.Context, key string, value []byte) error {
	conn, err := s.cache.GetContext(ctx)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	_, err = redis.DoWithTimeout(conn, cacheTimeout, "SET", key, value)
	if err != nil {
		log.Print(err)
	}
	return err
}