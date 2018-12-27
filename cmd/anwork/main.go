// This is the main anwork command line executable. This command line executable provides the
// ability to create, read, update, and delete anwork Task objects.
//
// Versioning is done with a single 32-bit integer. Version names start with a lowercase 'v' and are
// then followed by the number of the release. For example, the first version of the release was
// named _v1_. The second version of the release will be _v2_. There are no minor version
// numbers. This version number is controlled via the "version" global in the runner package. See
// the CLI command "anwork version" for more information.
package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/clock"
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

	l := log.New(ioutil.Discard, "ANWORK: ", log.Ldate|log.Ltime|log.Lshortfile)

	var repo task.Repo
	if address, ok := useApi(); ok {
		repo = client.New(l, address, wireAuth(l), wireCache(l))
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

func wireAuth(l *log.Logger) *auth.Client {
	privateKeyData, ok := os.LookupEnv("ANWORK_API_PRIVATE_KEY")
	if !ok {
		l.Fatal("must set ANWORK_API_PRIVATE_KEY")
	}

	block, _ := pem.Decode([]byte(privateKeyData))
	if expected := "RSA PRIVATE KEY"; block.Type != expected {
		l.Fatalf("unexpected PEM type: got %s, expected %s", block.Type, expected)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		l.Fatalf("failed to parse private key: %s", err.Error())
	}

	secret, ok := os.LookupEnv("ANWORK_API_SECRET")
	if !ok {
		l.Fatal("must set ANWORK_API_SECRET")
	}

	return auth.NewClient(clock.NewClock(), privateKey, []byte(secret))
}

func wireCache(l *log.Logger) *cache.Cache {
	var dir string
	if homeDir, ok := os.LookupEnv("HOME"); ok {
		dir = filepath.Join(homeDir, ".anwork")
		os.MkdirAll(dir, 0755)
	} else {
		dir = "."
	}

	return cache.New(filepath.Join(dir, "token-cache"))
}
