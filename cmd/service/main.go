// This is the ANWORK service. It runs an HTTP server and serves the ANWORK API.
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"reflect"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/auth"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/fs"
	"github.com/ankeesler/anwork/task/sql"
	_ "github.com/go-sql-driver/mysql"
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

	logger := lager.NewLogger("anwork-service")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.DEBUG))
	logger.Info("hey")

	repo := wireRepo(logger)

	clock := clock.NewClock()
	publicKey := getPublicKey(logger.Session("get-public-key"))
	secret := getSecret(logger.Session("get-secret"))
	authenticator := auth.NewServer(clock, rand.Reader, publicKey, secret)

	runner := http_server.New(
		address, api.New(logger.Session("api"), repo, authenticator))
	process := ifrit.Invoke(runner)
	logger.Info("running")

	logger.Fatal("process-exited", <-process.Wait())
}

func wireRepo(logger lager.Logger) task.Repo {
	var repo task.Repo
	if dsn, ok := os.LookupEnv("ANWORK_SQL_DSN"); ok {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			logger.Fatal("open-db-failure", err)
		}
		repo = sql.New(logger.Session("repo"), db)
	} else {
		repo = fs.New("/tmp/default-context")
	}
	return repo
}

func getPublicKey(logger lager.Logger) *rsa.PublicKey {
	publicKeyPEMBytes, ok := os.LookupEnv("ANWORK_API_PUBLIC_KEY")
	if !ok {
		msg := "could not read public key from ANWORK_API_PUBLIC_KEY env var"
		logger.Fatal("missing-env-var", errors.New(msg))
	}

	block, _ := pem.Decode([]byte(publicKeyPEMBytes))
	if block == nil {
		msg := "ANWORK_API_PUBLIC_KEY is in an invalid format"
		logger.Fatal("failed-to-decode-pem", errors.New(msg))
	}
	if expected := "PUBLIC KEY"; block.Type != expected {
		msg := fmt.Sprintf("unexpected PEM type: got %s, expected %s",
			block.Type, expected)
		logger.Fatal("invalid-pem-type", errors.New(msg))
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		msg := fmt.Sprintf("could not parse PKIX public key: %s", err.Error())
		logger.Fatal("failed-to-parse", errors.New(msg))
	}

	publicKey, ok := key.(*rsa.PublicKey)
	if !ok {
		msg := fmt.Sprintf("expected type *rsa.PublicKey from ANWORK_API_PUBLIC_KEY, got %s",
			reflect.TypeOf(key).String())
		logger.Fatal("wrong-key-type", errors.New(msg))
	}

	return publicKey
}

func getSecret(logger lager.Logger) []byte {
	secret, ok := os.LookupEnv("ANWORK_API_SECRET")
	if !ok {
		msg := "could not read secret from ANWORK_API_SECRET env var"
		logger.Fatal("missing-env-var", errors.New(msg))
	}

	return []byte(secret)
}
