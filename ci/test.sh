#!/bin/bash

set -ex

go vet ./...

ginkgo -race -r -cover
ANWORK_TEST_RUN_WITH_API=1 ginkgo -race integration
