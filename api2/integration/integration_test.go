package integration_test

import (
	"github.com/ankeesler/anwork/api2/client"
	"github.com/ankeesler/anwork/task2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	// TODO: stand up API.
	task2.RunRepoTests(func() task2.Repo {
		return client.New()
	})
})
