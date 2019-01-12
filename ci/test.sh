#!/bin/bash

set -ex

ginkgo -race -r -cover -skipPackage sql
ANWORK_TEST_RUN_WITH_API=1 ginkgo -race integration
