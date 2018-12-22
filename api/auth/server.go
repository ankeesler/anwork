package auth

import (
	"crypto/rsa"
	"fmt"
	"io"
	"time"

	"code.cloudfoundry.org/clock"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Server uses an RSA public key and a random byte-string secret to
// generate tokens and prove their authenticity. This Server generates
// tokens encrypted with the RSA public key and expects users to hold the
// matching RSA private key to decrypt the tokens. The Authenticate() method
// expects the token to be decrypted.
//
// This implementation uses JWT tokens according to RFC 7519. Then tokens are
// both encrypted and signed, according to RFC 7516 and 7515.
type Server struct {
	clock clock.Clock
	rand  io.Reader

	publicKey *rsa.PublicKey
	secret    []byte

	currentJTI string
}

// NewServer creates a new Server with a publicKey and a secret. It will use
// the provided clock to fill in the time-related claims of the JWT and the
// rand to populate the JWT ID (jti) field of the JWT.
func NewServer(
	clock clock.Clock,
	rand io.Reader,
	publicKey *rsa.PublicKey,
	secret []byte,
) *Server {
	return &Server{
		clock:     clock,
		rand:      rand,
		publicKey: publicKey,
		secret:    secret,
	}
}

func (s *Server) Authenticate(token string) error {
	parsed, err := jwt.ParseSigned(token)
	if err != nil {
		return fmt.Errorf("could not parse token: %s", err.Error())
	}

	claims := jwt.Claims{}
	if err := parsed.Claims(s.secret, &claims); err != nil {
		return fmt.Errorf("could not get claims: %s", err.Error())
	}

	if len(s.currentJTI) == 0 {
		return fmt.Errorf("null jti")
	}

	if err := claims.Validate(jwt.Expected{
		Issuer:  "anwork",
		Subject: "andrew",
		Time:    s.clock.Now(),
		ID:      s.currentJTI,
	}); err != nil {
		return fmt.Errorf("invalid claims: %s", err.Error())
	}

	return nil
}

func (s *Server) Token() (string, error) {
	// TODO: what signing algorithm should we be using?
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: s.secret}
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
			Key:       s.publicKey,
		},
		encrypterOptions,
	)
	if err != nil {
		return "", fmt.Errorf("could not create encrypter: %s", err.Error())
	}
	// TODO: how big should this be?
	r := make([]byte, 32)
	if n, err := s.rand.Read(r); n != len(r) {
		return "", fmt.Errorf("could not get %d random bytes: got %d", len(r), n)
	} else if err != nil {
		return "", fmt.Errorf("could not get %d random bytes: %s", len(r), err.Error())
	}

	now := s.clock.Now()
	s.currentJTI = string(r)
	claims := jwt.Claims{
		Issuer:    "anwork",
		Subject:   "andrew",
		Expiry:    jwt.NewNumericDate(now.Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        s.currentJTI,
	}
	token, err := jwt.SignedAndEncrypted(signer, enc).Claims(claims).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("could not sign and encrypt: %s", err.Error())
	}

	return token, nil
}
