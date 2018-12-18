// Package authenticator provides an implementation of an authentication mechanism
// for the ANWORK API.
package authenticator

import (
	"crypto/rsa"
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2/jwt"
)

// Authenticator uses an RSA private key and a random byte-string secret to
// generate tokens and prove their authenticity.
//
// This implementation uses JWT tokens according to RFC 7519. Then tokens are
// both encrypted and signed, according to RFC 7516 and 7515.
type Authenticator struct {
	privateKey *rsa.PrivateKey
	secret     []byte
}

// New creates a new Authenticator with a privateKey and a secret.
func New(privateKey *rsa.PrivateKey, secret []byte) *Authenticator {
	return &Authenticator{privateKey: privateKey, secret: secret}
}

func (a *Authenticator) Authenticate(token string) error {
	parsed, err := jwt.ParseSignedAndEncrypted(token)
	if err != nil {
		return fmt.Errorf("could not parse token: %s", err.Error())
	}

	nested, err := parsed.Decrypt(a.privateKey)
	if err != nil {
		return fmt.Errorf("could not decrypt token: %s", err.Error())
	}

	out := jwt.Claims{}
	if err := nested.Claims(a.secret, &out); err != nil {
		return fmt.Errorf("could not get claims: %s", err.Error())
	}

	if err := out.Validate(jwt.Expected{
		Issuer:  "anwork",
		Subject: "andrew",
		Time:    time.Now(),
	}); err != nil {
		return fmt.Errorf("invalid claims: %s", err.Error())
	}

	return nil
}

func (a *Authenticator) Token() (string, error) {
	return "", nil
}
