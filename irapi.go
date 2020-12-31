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
	"os"
	"strings"
)

// UserAgent is the value given for the User-Agent
//const UserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0"
const UserAgent = "irapi/1.0 +https://github.com/LeoAdamek/irapi"

var (

	// ErrLoginFailed is a generic error for when login fails due to credentials or system failure
	ErrLoginFailed = errors.New("failed login")

	// ErrMaintenance is an error returned when iRacing is down for maintanenace
	ErrMaintenance = errors.New("iRacing is offline for maintaneance")

	// ErrTooManyRequests is an error returned when iRacing rejects requests due to volume
	ErrTooManyRequests = errors.New("too many requests")
)

// IRacing is an instance of an API client for the iRacing Service
type IRacing struct {
	http                *http.Client
	credentialsProvider CredentialsProvider
	BeforeFuncs         []BeforeFunc
	AfterFuncs          []AfterFunc
}

// BeforeFunc is a function which runs before a request is sent
//
// These can be used with `IRacing.BeforeRequest()` to add middleware before a request is made.
// BeforeFunc handlers are responsible for preserving the content of `req.Body` if they consume it.
type BeforeFunc func(ctx context.Context, req *http.Request) error

// AfterFunc is a function which is fun after a response is received
//
// These can be used with `IRacing.AfterResponse()` to add middleware after a response is received.
// AfterFunc handlers are responsible for preserving the content of `res.Body` if they consume it.
type AfterFunc func(ctx context.Context, req *http.Request, res *http.Response) error

// Host is the address where the iRacing service is hosted
const Host = "https://members.iracing.com"

// New crates a new iRacing API client instance
func New(credentials CredentialsProvider) *IRacing {

	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Jar: jar,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
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
	req.Header.Set("Origin", "members.iracing.com")
	req.Header.Set("Referer", Host+"/membersite/login.jsp")

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
		if res.StatusCode == http.StatusTooManyRequests {
			return res, ErrTooManyRequests
		}

		if res.StatusCode >= 500 {
			err = errors.New("error server response")
		}
	}

	if res.Header.Get("X-Maintenance-Mode") == "true" {
		return res, ErrMaintenance
	}

	return res, err
}

func (c *IRacing) json(ctx context.Context, method, path string, body, into interface{}) error {
	var reader io.Reader

	if body != nil {
		switch b := body.(type) {
		case io.Reader:
			reader = b
		default:
			buffer := new(bytes.Buffer)
			if err := json.NewEncoder(buffer).Encode(body); err != nil {
				return err
			}

			reader = buffer
		}
	}

	req, _ := http.NewRequestWithContext(ctx, method, Host+path, reader)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := c.do(ctx, req)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode >= 500 {
			return errors.New("server returned an error")
		}

		return errors.New("server rejected our request")
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
//
// NOTE: Middleware is not invoked for Login requests.
func (c *IRacing) Login(ctx context.Context) error {

	credentials, err := c.credentialsProvider()

	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("username", credentials.Username)
	params.Set("password", credentials.Password)
	params.Set("utcoffset", "0")
	params.Set("todaysdate", "")

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, Host+"/membersite/Login", strings.NewReader(params.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.do(ctx, req)

	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusOK {
		return nil
	} else if res.StatusCode == http.StatusFound {
		redirect := res.Header.Get("Location")

		if redirect == "https://members.iracing.com/membersite/failedlogin.jsp" {

			io.Copy(os.Stderr, res.Body)

			return ErrLoginFailed
		}
		return nil
	}

	if res.Header.Get("X-Maintenance-Mode") == "true" {
		return ErrMaintenance
	}

	return fmt.Errorf("expected to get an OK response to login, instead got %s", res.Status)
}
