package URLC

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Service struct {
	MLURL string
}
type Req struct {
	Url string
}
type Resp struct {
	Verdict string
}
func NewService(url string) *Service{
	return &Service{url}
}

func (s*Service) CheckURL(url string) (string, error){
	log.Println("Enter chec URL")
	data := Req{Url: url}
	request, err := json.Marshal(data)
	if err != nil{
		return "", err
	}
	response, err := http.Post(s.MLURL, "application/json", bytes.NewBuffer(request))
	if err != nil{
		return "", err
	}
	var strResponse Resp
	err = json.NewDecoder(response.Body).Decode(&strResponse)
	if err != nil{
		return "", err
	}
	return strResponse.Verdict, nil
}