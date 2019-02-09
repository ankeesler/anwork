// Package readtoken contains a utility program for reading the encrypted JWT
// token from the token cache. It requires that the ANWORK_API_SECRET and
// ANWORK_API_PRIVATE_KEY environmental variables are set properly.
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func main() {
	privateKey, err := readPrivateKey()
	if err != nil {
		die(fmt.Sprintf("read privatekey: %s", err.Error()))
	}

	secret, ok := os.LookupEnv("ANWORK_API_SECRET")
	if !ok {
		die("ANWORK_API_SECRET must be set")
	}

	encryptedToken := getTokenCacheFileData()
	decryptedToken := decryptToken(encryptedToken, privateKey, []byte(secret))
	printDecryptedToken(decryptedToken)
}

func die(msg string) {
	fmt.Println("error:", msg)
	os.Exit(1)
}

func getTokenCacheFileData() string {
	var dir string
	if homeDir, ok := os.LookupEnv("HOME"); ok {
		dir = filepath.Join(homeDir, ".anwork")
		os.MkdirAll(dir, 0755)
	} else {
		dir = "."
	}

	file := filepath.Join(dir, "token-cache")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		die(fmt.Sprintf("read token cache file: %s", err.Error()))
	}

	return string(data)
}

func decryptToken(
	encryptedToken string,
	privateKey *rsa.PrivateKey,
	secret []byte,
) string {
	parsed, err := jwt.ParseSignedAndEncrypted(encryptedToken)
	if err != nil {
		die(fmt.Sprintf("could not parse token: %s", err.Error()))
	}

	nested, err := parsed.Decrypt(privateKey)
	if err != nil {
		die(fmt.Sprintf("could not decrypt token: %s", err.Error()))
	}

	claims := jwt.Claims{}
	if err := nested.Claims(secret, &claims); err != nil {
		die(fmt.Sprintf("could not verify claims: %s", err.Error()))
	}

	signingKey := jose.SigningKey{Algorithm: jose.HS512, Key: secret}
	signerOptions := (&jose.SignerOptions{}).WithType("JWT")
	signer, err := jose.NewSigner(signingKey, signerOptions)
	if err != nil {
		die(fmt.Sprintf("could not make signer: %s", err.Error()))
	}

	decryptedToken, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		die(fmt.Sprintf("signing token: %s", err.Error()))
	}

	return decryptedToken
}

func printDecryptedToken(token string) {
	fmt.Println("token:", token)

	pieces := strings.Split(token, ".")
	if len(pieces) != 3 {
		die(fmt.Sprintf("split: expected 3 pieces, got %d", len(pieces)))
	}

	body, err := base64.RawURLEncoding.DecodeString(pieces[1])
	if err != nil {
		die(fmt.Sprintf("decoding body failed: %s", err.Error()))
	}

	data := struct {
		Exp int64  `json:"exp"`
		Iat int64  `json:"iat"`
		Iss string `json:"iss"`
		Jti string `json:"jti"`
		Nbf int64  `json:"nbf"`
		Sub string `json:"sub"`
	}{}
	if err := json.Unmarshal(body, &data); err != nil {
		die(fmt.Sprintf("unmarshal jwt body: %s", err.Error()))
	}

	fmt.Println("body:")
	fmt.Printf("  iat: %d (%s)\n", data.Iat, time.Unix(data.Iat, 0).String())
	fmt.Printf("  nbf: %d (%s)\n", data.Nbf, time.Unix(data.Nbf, 0).String())
	fmt.Printf("  exp: %d (%s)\n", data.Exp, time.Unix(data.Exp, 0).String())
	fmt.Printf("  iss: %s\n", data.Iss)
	fmt.Printf("  sub: %s\n", data.Sub)
	fmt.Printf("  jti: %s\n", data.Jti)
}

func readPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyData, ok := os.LookupEnv("ANWORK_API_PRIVATE_KEY")
	if !ok {
		return nil, errors.New("must set ANWORK_API_PRIVATE_KEY")
	}

	block, _ := pem.Decode([]byte(privateKeyData))
	if block == nil {
		return nil, errors.New("failed to decode private key PEM data")
	}
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		return nil, fmt.Errorf("unexpected PEM type: got %s, expected %s",
			block.Type, expected)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %s", err.Error())
	}

	return privateKey, nil
}
