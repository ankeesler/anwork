// Package authenticator provides an implementation of an authentication mechanism
// for the ANWORK API.
package authenticator

import (
	"crypto/rsa"
	"fmt"
	"io"
	"time"

	"code.cloudfoundry.org/clock"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Authenticator uses an RSA public key and a random byte-string secret to
// generate tokens and prove their authenticity. This Authenticator generates
// tokens encrypted with the RSA public key and expects users to hold the
// matching RSA private key to decrypt the tokens. The Authenticate() method
// expects the token to be decrypted.
//
// This implementation uses JWT tokens according to RFC 7519. Then tokens are
// both encrypted and signed, according to RFC 7516 and 7515.
type Authenticator struct {
	clock clock.Clock
	rand  io.Reader

	publicKey *rsa.PublicKey
	secret    []byte
}

// New creates a new Authenticator with a publicKey and a secret. It will use
// the provided rand to populate the JWT ID (jti) field of the JWT.
func New(clock clock.Clock, rand io.Reader, publicKey *rsa.PublicKey, secret []byte) *Authenticator {
	return &Authenticator{
		clock:     clock,
		rand:      rand,
		publicKey: publicKey,
		secret:    secret,
	}
}

func (a *Authenticator) Authenticate(token string) error {
	parsed, err := jwt.ParseSigned(token)
	if err != nil {
		return fmt.Errorf("could not parse token: %s", err.Error())
	}

	claims := jwt.Claims{}
	if err := parsed.Claims(a.secret, &claims); err != nil {
		return fmt.Errorf("could not get claims: %s", err.Error())
	}

	if err := claims.Validate(jwt.Expected{
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
		return "", fmt.Errorf("could not create signer: %s", err.Error())
	}

	// TODO: what encryption algorithm should we be using?
	encrypterOptions := (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT")
	enc, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       a.publicKey,
		},
		encrypterOptions,
	)
	if err != nil {
		return "", fmt.Errorf("could not create encrypter: %s", err.Error())
	}

	// TODO: how big should this be?
	r := make([]byte, 32)
	if n, err := a.rand.Read(r); n != 32 {
		return "", fmt.Errorf("could not get 32 random bytes: got %d", n)
	} else if err != nil {
		return "", fmt.Errorf("could not get 32 random bytes: %s", err.Error())
	}

	now := a.clock.Now()
	claims := jwt.Claims{
		Issuer:    "anwork",
		Subject:   "andrew",
		Expiry:    jwt.NewNumericDate(now.Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token, err := jwt.SignedAndEncrypted(signer, enc).Claims(claims).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("could not sign and encrypt: %s", err.Error())
	}

	return token, nil
}
