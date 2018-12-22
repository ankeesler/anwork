package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"code.cloudfoundry.org/clock"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Client provides client-side authentication functionality.
//
// It has the ability that a JWT has been generated by a Server's
// Token() function. That is, the Client's Validate() function will
// decrypt a JWT, cryptographically verify it, and then sign it again
// with the same secret in order to be sent back to the Server.
type Client struct {
	clock clock.Clock

	privateKey *rsa.PrivateKey
	secret     []byte
}

// New creates a new Client.
func NewClient(clock clock.Clock, privateKey *rsa.PrivateKey, secret []byte) *Client {
	return &Client{clock: clock, privateKey: privateKey, secret: secret}
}

func (c *Client) Validate(token string) (string, error) {
	parsed, err := jwt.ParseSignedAndEncrypted(token)
	if err != nil {
		return "", fmt.Errorf("could not parse token: %s", err.Error())
	}

	nested, err := parsed.Decrypt(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decrypt token: %s", err.Error())
	}

	claims := jwt.Claims{}
	if err := nested.Claims(c.secret, &claims); err != nil {
		return "", fmt.Errorf("could not verify claims: %s", err.Error())
	}

	if err := claims.Validate(jwt.Expected{
		Issuer:  "anwork",
		Subject: "andrew",
		Time:    time.Now(),
	}); err != nil {
		return "", fmt.Errorf("invalid claims: %s", err.Error())
	}

	// TODO: what signing algorithm should we be using?
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: c.secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	if err != nil {
		return "", fmt.Errorf("could not create signer: %s", err.Error())
	}

	return jwt.Signed(signer).Claims(claims).CompactSerialize()
}
