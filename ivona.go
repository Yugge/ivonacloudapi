package ivonacloudapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gopkg.in/amz.v1/aws"
	"io/ioutil"
	"net/http"
	"time"
)

var ENDPOINT_EU_WEST Endpoint = Endpoint{"tts.eu-west-1.ivonacloud.com", "eu-west-1"}
var ENDPOINT_US_EAST Endpoint = Endpoint{"tts.us-east-1.ivonacloud.com", "us-east-1"}

type Endpoint struct {
	URI  string
	Name string
}

type IvonaClient struct {
	Keys     ApiKeys
	Endpoint Endpoint
}

type ApiKeys struct {
	Access string
	Secret string
}

type listVoicesRequest struct {
	Voice map[string]string
}

type createSpeechRequest struct {
	Input map[string]string
	Voice map[string]string
}

type request struct {
	Request string
}
type param struct {
	Key   string
	Value string
}
type Voice struct {
	Name     string
	Language string
	Gender   string
}

func NewIvonaClient(access, secret string, endpoint Endpoint) *IvonaClient {
	return &IvonaClient{ApiKeys{access, secret}, endpoint}
}

func (i *IvonaClient) ListVoices(name, language, gender string) {
	voice := make(map[string]string)
	if name != "" {
		voice["Name"] = name
	}
	if language != "" {
		voice["Language"] = language
	}
	if gender != "" {
		voice["Gender"] = gender
	}
	l := listVoicesRequest{voice}
	js, _ := json.Marshal(l)
	fmt.Println(string(js))
	resp := i.makeRequest(js, "ListVoices")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func (i *IvonaClient) CreateSpeech(text string, v Voice) []byte {
	voice := make(map[string]string)
	if v.Name != "" {
		voice["Name"] = v.Name
	}
	if v.Language != "" {
		voice["Language"] = v.Language
	}
	if v.Gender != "" {
		voice["Gender"] = v.Gender
	}
	input := make(map[string]string)
	input["Data"] = text
	l := createSpeechRequest{input, voice}
	js, _ := json.Marshal(l)
	fmt.Println(string(js))
	resp := i.makeRequest(js, "CreateSpeech")
	body, _ := ioutil.ReadAll(resp.Body)
	return body

}

func (i *IvonaClient) makeRequest(payload []byte, request string) *http.Response {
	req, _ := http.NewRequest("POST", "https://"+i.Endpoint.URI+"/"+request, bytes.NewReader(payload))
	hasher := sha256.New()
	hasher.Write(payload)
	req.Header.Add("x-amz-content-sha256", fmt.Sprintf("%x", hasher.Sum(nil)))
	req.Header.Add("x-amz-date", time.Now().UTC().Format("20060102T150405Z"))
	req.Header.Add("host", i.Endpoint.URI)
	err := aws.SignV4(req, aws.Auth{i.Keys.Access, i.Keys.Secret}, i.Endpoint.Name)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}
