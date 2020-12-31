package irapi

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// Credentials represents the credentials required to access the iRacing service
type Credentials struct {
	Username string
	Password string
}

// CredentialsProvider is a function which attempts to return a set of Credentials
// This allows credentials to be seamlessly provided by any source, with providers for
// static credentials, text-file credntials and environment variable credentials included.
// Other providers can be easily created for sources like Vault, AWS Secrets etc.
type CredentialsProvider func() (*Credentials, error)

// StaticCredentialsProvider creates a new credential provider returning the given credentials
func StaticCredentialsProvider(username, password string) CredentialsProvider {
	return func() (*Credentials, error) {
		return &Credentials{
			Username: username,
			Password: password,
		}, nil
	}
}

// FileCredentialsProvider provides credentials from a file containing the credentials
// This is best used in conjunction with a secrets manager like Kubernetes.
//
// File Format:	`{{.Username}},{{.Password}}`
//
func FileCredentialsProvider(path string) CredentialsProvider {
	return func() (*Credentials, error) {
		data, err := ioutil.ReadFile(path)

		if err != nil {
			return nil, err
		}

		parts := strings.SplitN(string(data), ",", 2)

		return &Credentials{
			Username: parts[0],
			Password: parts[1],
		}, nil
	}
}

// EnvironmentCredentialsProvider gets credentials from the environment
// using the environment variables `IRACING_USERNAME` and `IRACING_PASSWORD`
func EnvironmentCredentialsProvider() (*Credentials, error) {
	username := os.Getenv("IRACING_USERNAME")
	password := os.Getenv("IRACING_PASSWORD")

	if username == "" {
		return nil, errors.New("no username found in $IRACING_USERNAME")
	}

	if password == "" {
		return nil, errors.New("no password found in $IRACING_PASSWORD")
	}

	return &Credentials{
		Username: username,
		Password: password,
	}, nil
}
