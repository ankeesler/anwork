// Package authenticator provides an implementation of an authentication mechanism
// for the ANWORK API.
package authenticator

import (
	"crypto/rsa"
	"fmt"
	"time"

	"code.cloudfoundry.org/clock"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Authenticator uses an RSA private key and a random byte-string secret to
// generate tokens and prove their authenticity.
//
// This implementation uses JWT tokens according to RFC 7519. Then tokens are
// both encrypted and signed, according to RFC 7516 and 7515.
type Authenticator struct {
	clock clock.Clock

	privateKey *rsa.PrivateKey
	secret     []byte
}

// New creates a new Authenticator with a privateKey and a secret.
func New(clock clock.Clock, privateKey *rsa.PrivateKey, secret []byte) *Authenticator {
	return &Authenticator{
		clock:      clock,
		privateKey: privateKey,
		secret:     secret,
	}
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
		Time:    a.clock.Now(),
	}); err != nil {
		return fmt.Errorf("invalid claims: %s", err.Error())
	}

	return nil
}

func (a *Authenticator) Token() (string, error) {
	// TODO: what signing algorithm should we be using?
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: a.secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	if err != nil {
		return "", err
	}

	// TODO: what encryption algorithm should we be using?
	encrypterOptions := (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT")
	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       &a.privateKey.PublicKey,
		},
		encrypterOptions,
	)
	if err != nil {
		return "", err
	}

	now := a.clock.Now()
	claims := jwt.Claims{
		Issuer:    "anwork",
		Subject:   "andrew",
		Expiry:    jwt.NewNumericDate(now.Add(time.Second * 1)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token, err := jwt.SignedAndEncrypted(signer, enc).Claims(claims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}
