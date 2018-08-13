#!/bin/bash

set -ex

ginkgo -r -cover
ANWORK_TEST_RUN_WITH_API=1 ginkgo integration
