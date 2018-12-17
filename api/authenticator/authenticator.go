// Package authenticator provides an implementation of an authentication mechanism
// for the ANWORK API.
package authenticator

import "net/http"

type Authenticator struct {
}

func New() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) Authenticate(r *http.Request) error {
	return nil
}

func (a *Authenticator) Token(r *http.Request) (string, error) {
	return "", nil
}
