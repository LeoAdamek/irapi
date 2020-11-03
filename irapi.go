package irapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// UserAgent is the value given for the User-Agent
const UserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0"

// IRacing is an instance of an API client for the iRacing Service
type IRacing struct {
	http                *http.Client
	credentialsProvider CredentialsProvider
	BeforeFuncs         []BeforeFunc
	AfterFuncs          []AfterFunc
}

// BeforeFunc is a function which runs before a request is sent
type BeforeFunc func(ctx context.Context, req *http.Request) error

// AfterFunc is a function which is fun after a response is received
type AfterFunc func(ctx context.Context, req *http.Request, res *http.Response) error

// Host is the address where the iRacing service is hosted
const Host = "https://members.iracing.com"

// New crates a new iRApi instance
func New(ctx context.Context, credentials CredentialsProvider) *IRacing {

	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: jar,
	}

	return &IRacing{
		http:                client,
		credentialsProvider: credentials,
	}
}

// BeforeRequest adds a new BeforeFunc to the chain
func (c *IRacing) BeforeRequest(f BeforeFunc) {
	c.BeforeFuncs = append(c.BeforeFuncs, f)
}

// AfterResponse adds a new AfterFunc to the response chain
func (c *IRacing) AfterResponse(f AfterFunc) {
	c.AfterFuncs = append(c.AfterFuncs, f)
}

// SetHTTP overrides the HTTP client for the API instance
func (c *IRacing) SetHTTP(client *http.Client) {
	c.http = client
}

// do run an HTTP Request
func (c IRacing) do(ctx context.Context, req *http.Request) (*http.Response, error) {

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Cache-Control", "no-cache")

	for _, f := range c.BeforeFuncs {
		if err := f(ctx, req); err != nil {
			return nil, err
		}
	}

	res, err := c.http.Do(req)

	for _, f := range c.AfterFuncs {
		if err := f(ctx, req, res); err != nil {
			return nil, err
		}
	}

	if res.StatusCode >= 400 {
		err = errors.New("error server response")
	}

	return res, err
}

func (c *IRacing) json(ctx context.Context, method, path string, body, into interface{}) error {
	var reader io.Reader

	if body != nil {
		buffer := new(bytes.Buffer)
		if err := json.NewEncoder(buffer).Encode(body); err != nil {
			return err
		}

		reader = buffer
	}

	req, _ := http.NewRequestWithContext(ctx, method, Host+path, reader)

	req.Header.Set("Accept", "application/json")

	res, err := c.do(ctx, req)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New("server returned an error")
	}

	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	decoded, err := url.QueryUnescape(string(content))

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(decoded), into)
}

// Login will log into the iRacing Service
func (c *IRacing) Login(ctx context.Context) error {

	credentials, err := c.credentialsProvider()

	if err != nil {
		return err
	}

	params := make(url.Values)
	params.Set("username", credentials.Username)
	params.Set("password", credentials.Password)
	params.Set("AUTOLOGIN", "on")
	params.Set("utcoffset", "0")
	params.Set("todaysdate", "")

	body := params.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, Host+"/membersite/Login", strings.NewReader(body))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", Host+"/membersite/login.jsp")

	res, err := c.do(ctx, req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected to get an OK response to login, instead got %s", res.Status)
	}

	return nil
}
