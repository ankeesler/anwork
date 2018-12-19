package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/authenticator"
	"github.com/ankeesler/anwork/lag"
	"github.com/ankeesler/anwork/task/fs"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

func main() {
	var address string
	if port, ok := os.LookupEnv("PORT"); ok {
		address = fmt.Sprintf(":%s", port)
	} else {
		address = ":12345"
	}

	log := log.New(os.Stdout, "ANWORK Service: ", log.Ldate|log.Ltime)
	l := lag.New(log, lag.D)
	l.P(lag.I, "hey")

	repo := fs.New("/tmp/default-context")

	// clock := clock.NewClock()
	// publicKey := getPublicKey(log)
	// secret := getSecret(log)
	// authenticator := authenticator.New(clock, rand.Reader, publicKey, secret)
	authenticator := authenticator.NullAuthenticator{}

	runner := http_server.New(address, api.New(l, repo, authenticator))
	process := ifrit.Invoke(runner)
	l.P(lag.I, "running")

	log.Fatal(<-process.Wait())
}

func getPublicKey(log *log.Logger) *rsa.PublicKey {
	publicKeyPEMBytes, ok := os.LookupEnv("ANWORK_API_PUBLIC_KEY")
	if !ok {
		log.Fatalf("could not read public key file from ANWORK_API_PUBLIC_KEY env var")
	}

	block, _ := pem.Decode([]byte(publicKeyPEMBytes))
	if expected := "PUBLIC KEY"; block.Type != expected {
		log.Fatalf("unexpected PEM type: got %s, expected %s", block.Type, expected)
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("could not parse PKCS#1 public key: %s", err.Error())
	}

	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("expected type *rsa.PublicKey from ANWORK_API_PUBLIC_KEY, got %s",
			reflect.TypeOf(key).String())
	}

	return publicKey
}

func getSecret(log *log.Logger) []byte {
	secret, ok := os.LookupEnv("ANWORK_API_SECRET")
	if !ok {
		log.Fatalf("could not read secret from ANWORK_API_SECRET env var")
	}

	return []byte(secret)
}
