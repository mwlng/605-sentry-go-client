package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/anaskhan96/soup"
)

const (
	authUri   = "/oauth/authorize"
	signinUri = "/users/sign_in"
)

type SentryClient struct {
	client    http.Client
	serverUrl string
}

func New() *SentryClient {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)

	}
	httpClient := http.Client{
		Jar: jar,
	}
	return &SentryClient{
		client: httpClient,
	}
}

func (s *SentryClient) Login(serverUrl, userName, password string) {
	req, _ := http.NewRequest("GET", serverUrl, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	authenticityToken := extract_authenticity_token(bodyBytes)

	body := url.Values{
		"authenticity_token": {authenticityToken},
		"user[email]":        {userName},
		"user[password]":     {password},
		"user[remember_me]":  {"0"},
		"commit":             {"Login"},
	}
	signinUrl := serverUrl + signinUri
	req, _ = http.NewRequest("POST", signinUrl, strings.NewReader(body.Encode()))

	resp, err = s.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Login failed!")
	}
	s.serverUrl = serverUrl
}

func extract_authenticity_token(bodyBytes []byte) string {
	doc := soup.HTMLParse(string(bodyBytes))
	attrs := doc.Find("input", "name", "authenticity_token").Attrs()
	return attrs["value"]
}

func (s *SentryClient) GetAccessToken(clientID, responseType, redirectUrl string) []string {
	authorizeUrl := s.serverUrl + authUri
	reqUrl := fmt.Sprintf(
		"%s?client_id=%s&response_type=%s&redirect_uri=%s&scopes=public&state=xyz",
		authorizeUrl,
		clientID,
		responseType)

	req, err := http.NewRequest("GET", reqUrl, nil)

	s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		if resp.StatusCode == http.StatusFound {
			return resp.Header["Location"]
		} else {
			log.Fatal("Can't get access token")
		}
	}
	return []string{}
}
