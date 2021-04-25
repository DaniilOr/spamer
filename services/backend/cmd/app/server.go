package app

import (
	"context"
	"encoding/json"
	"github.com/DaniilOr/spamer/services/backend/pkg/auth"
	"github.com/DaniilOr/spamer/services/backend/pkg/classifier"
	"github.com/DaniilOr/spamer/services/backend/pkg/spam"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"io/ioutil"
	"log"
	"net/http"
)
type Params struct {
	Interval    int64 	`json:"interval"`
	Url string `json:"url"`
	NumStreams int64 `json:"num_streams"`
}
type Resp struct {
	Result string `json:"result"`
}
type Req struct {
	Url string
}
type Server struct {
	authSvc         *auth.Service
	classifierSvc   *classifier.Service
	spamerSvc       *spam.Service
	mux             chi.Router
}
type UserCreds struct{
	Login    string `json:"login"`
	Password string `json:"password"`
}
func NewServer(authSvc *auth.Service, classifierSvc *classifier.Service, spamerSvc *spam.Service, mux chi.Router) *Server {
	return &Server{authSvc: authSvc, classifierSvc: classifierSvc, spamerSvc:spamerSvc, mux: mux}
}

func (s *Server) Init() error {
	s.mux.Use(middleware.Logger)
	s.mux.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Post("/token", s.token)
		r.With(Auth(func(ctx context.Context, token string) (int64, error) {
			return s.authSvc.Id(ctx, token)
		})).Post("/spam", s.spam)
	})
	s.mux.Route("/classify", func(r chi.Router){
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins:   []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Post("/url", s.classifyURL)
		r.Post("/sms", s.classifySMS)
	})

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) token(writer http.ResponseWriter, request *http.Request) {
	var uCreds UserCreds
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&uCreds)
	if err != nil {
		log.Print("can't parse form")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("%v", uCreds)

	token, err := s.authSvc.Token(request.Context(), uCreds.Login, uCreds.Password)
	if err != nil {
		log.Printf("Auth Service returns error: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data := &tokenDTO{Token: token, Exp: 600, Login: "user"}
	respBody, err := json.Marshal(data)
	if err != nil {
		log.Printf("can't marshall data: %v", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(respBody)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) spam(writer http.ResponseWriter, request *http.Request) {
	_, err := AuthFrom(request.Context())
	if err != nil {
		log.Printf("can't find userID in context: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	var params Params
	err = json.Unmarshal(data, &params)
	if err != nil{
		log.Printf("fail to unmarshal spam params: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := s.spamerSvc.Spam(request.Context(), params.Url, params.Interval, params.NumStreams)
	if err != nil {
		log.Printf("Spamer Service returns error: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err = json.Marshal(Resp{Result: result})
	if err != nil{
		log.Printf("Cannot marhsal spamer result: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) classifyURL(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	var r Req
	err = json.Unmarshal(data, &r)
	if err != nil{
		log.Printf("fail to unmarshal url params: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := s.classifierSvc.CheckURL(request.Context(), r.Url)
	if err != nil {
		log.Printf("Classifier URl Service returns error: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err = json.Marshal(Resp{Result: result})
	if err != nil{
		log.Printf("Cannot marhsal classifier result: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}
func (s *Server) classifySMS(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	var r Req
	err = json.Unmarshal(data, &r)
	if err != nil{
		log.Printf("fail to unmarshal url params: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := s.classifierSvc.CheckSMS(request.Context(), r.Url)
	if err != nil {
		log.Printf("Classifier SMA Service returns error: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err = json.Marshal(Resp{Result: result})
	if err != nil{
		log.Printf("Cannot marhsal classifier result: %v", err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
}
