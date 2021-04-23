package spam

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Service struct {
	PYURL string
}
type Params struct {
	Interval    int64 	`json:"interval"`
	Url string `json:"url"`
	NumStreams int64 `json:"num_streams"`
}
type Resp struct {
	Result string `json:"result"`
}
func NewService(url string) *Service{
	return &Service{PYURL: url}
}
func (s*Service) Spam(url string, interval int64, numStreams int64) (string, error){
	data := Params{Interval: interval, Url: url, NumStreams: numStreams}
	request, err := json.Marshal(data)
	if err != nil{
		return "", err
	}
	response, err := http.Post(s.PYURL, "application/json", bytes.NewBuffer(request))
	if err != nil{
		return "", err
	}
	var strResponse Resp
	err = json.NewDecoder(response.Body).Decode(&strResponse)
	if err != nil{
		return "", err
	}
	return strResponse.Result, nil
}