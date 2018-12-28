// This is the main ANWORK command line executable. This command line executable provides the
// ability to create, read, update, and delete anwork Task objects.
package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/ankeesler/anwork/api/auth"
	"github.com/ankeesler/anwork/api/client"
	"github.com/ankeesler/anwork/api/client/cache"
	"github.com/ankeesler/anwork/manager"
	runner "github.com/ankeesler/anwork/runner"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/fs"
)

var (
	buildHash = "(dev)"
	buildDate = "???"
)

type rootFlagValue struct {
	value string
}

func (rfv *rootFlagValue) String() string {
	if len(rfv.value) == 0 {
		if homeDir, ok := os.LookupEnv("HOME"); ok {
			dir := filepath.Join(homeDir, ".anwork")
			os.MkdirAll(dir, 0755)
			return dir
		} else {
			return "."
		}
	}
	return rfv.value
}

func (rfv *rootFlagValue) Set(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("Cannot have a root flag with length 0!")
	}
	rfv.value = value
	return nil
}

type debugWriter struct {
	debug bool
}

func (dw *debugWriter) Write(data []byte) (int, error) {
	if dw.debug {
		return fmt.Print(string(data))
	}
	return 0, nil
}

func main() {
	var (
		context string
		root    rootFlagValue
		dw      debugWriter
	)

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flags.BoolVar(&dw.debug, "d", false, "Enable debug printing")

	flags.StringVar(&context, "c", "default-context", "Set the persistence context")
	flags.Var(&root, "o", "Set the persistence root directory")

	flags.Usage = func() {
		fmt.Println("Usage of anwork")
		fmt.Println("Flags")
		flags.SetOutput(os.Stdout)
		flags.PrintDefaults()
		fmt.Println("Commands")
		runner.Usage(os.Stdout)
	}

	if err := flags.Parse(os.Args[1:]); err == flag.ErrHelp {
		// Looks like help is printed by the flag package...
		os.Exit(0)
	} else if err != nil {
		// I think the flag package prints out the error and the usage...
		os.Exit(1)
	}

	if flags.NArg() == 0 {
		// If there are no arguments, return success. People might use this to simply check if the anwork
		// executable is on their machine.
		flags.Usage()
		os.Exit(0)
	}

	var logLevel lager.LogLevel
	if dw.debug {
		logLevel = lager.DEBUG
	} else {
		logLevel = lager.FATAL
	}

	logger := lager.NewLogger("anwork")
	logger.RegisterSink(lager.NewPrettySink(os.Stdout, logLevel))

	var repo task.Repo
	if address, ok := useApi(); ok {
		repo = client.New(
			logger.Session("api-client"),
			address,
			wireAuth(logger.Session("wire-auth")),
			wireCache(logger.Session("wire-cache")),
		)
	} else {
		repo = fs.New(filepath.Join(root.String(), context))
	}

	clock := clock.NewClock()
	m := manager.New(repo, clock)

	r := runner.New(&runner.BuildInfo{Hash: buildHash, Date: buildDate}, m, os.Stdout, &dw)
	if err := r.Run(flags.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func useApi() (string, bool) {
	return os.LookupEnv("ANWORK_API_ADDRESS")
}

func wireAuth(logger lager.Logger) *auth.Client {
	privateKeyData, ok := os.LookupEnv("ANWORK_API_PRIVATE_KEY")
	if !ok {
		msg := "must set ANWORK_API_PRIVATE_KEY"
		logger.Fatal("missing-private-key-env-var", errors.New(msg))
	}

	block, _ := pem.Decode([]byte(privateKeyData))
	if block == nil {
		msg := "failed to decode private key PEM data"
		logger.Fatal("failed-to-decode-private-key-pem", errors.New(msg))
	}
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		msg := fmt.Sprintf("unexpected PEM type: got %s, expected %s",
			block.Type, expected)
		logger.Fatal("invalid-private-key-pem-type", errors.New(msg))
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Fatal("failed-to-parse-private-key", err)
	}

	secret, ok := os.LookupEnv("ANWORK_API_SECRET")
	if !ok {
		msg := "must set ANWORK_API_SECRET"
		logger.Fatal("missing-secret-env-var", errors.New(msg))
	}

	return auth.NewClient(clock.NewClock(), privateKey, []byte(secret))
}

func wireCache(logger lager.Logger) *cache.Cache {
	var dir string
	if homeDir, ok := os.LookupEnv("HOME"); ok {
		dir = filepath.Join(homeDir, ".anwork")
		os.MkdirAll(dir, 0755)
	} else {
		dir = "."
	}

	return cache.New(filepath.Join(dir, "token-cache"))
}
