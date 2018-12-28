package auth

import (
	"crypto/rsa"
	"fmt"

	jose "gopkg.in/square/go-jose.v2"
)

func signer(secret []byte) (jose.Signer, error) {
	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %s", err.Error())
	}

	return signer, nil
}

func encrypter(publicKey *rsa.PublicKey) (jose.Encrypter, error) {
	encrypterOptions := (&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT")
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP_256,
			Key:       publicKey,
		},
		encrypterOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create encrypter: %s", err.Error())
	}

	return encrypter, nil
}
