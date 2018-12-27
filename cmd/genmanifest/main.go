package main

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type manifest struct {
	Applications []application `yaml:"applications"`
}

type application struct {
	Buildpacks              []string          `yaml:"buildpacks"`
	Memory                  string            `yaml:"memory"`
	DiskQuota               string            `yaml:"disk_quota"`
	Command                 string            `yaml:"command"`
	HealthCheckType         string            `yaml:"health-check-type"`
	HealthCheckHTTPEndpoint string            `yaml:"health-check-http-endpoint"`
	Environment             map[string]string `yaml:"env"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <public-key> <secret>", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	publicKey := os.Args[1]
	secret := os.Args[2]

	manifest := makeManifest(publicKey, secret)
	bytes, err := yaml.Marshal(manifest)
	if err != nil {
		fmt.Println("error: could not marshal manifest")
		fmt.Println("  err:", err)
		fmt.Println("  manifest:", manifest)
	}

	fmt.Println("---")
	fmt.Println(string(bytes))
}

func makeManifest(publicKey, secret string) *manifest {
	return &manifest{
		Applications: []application{
			application{
				Buildpacks:              []string{"binary_buildpack"},
				Memory:                  "256M",
				DiskQuota:               "256M",
				Command:                 "./main",
				HealthCheckType:         "http",
				HealthCheckHTTPEndpoint: "/api/v1/health",
				Environment: map[string]string{
					"ANWORK_API_PUBLIC_KEY": publicKey,
					"ANWORK_API_SECRET":     secret,
				},
			},
		},
	}
}
