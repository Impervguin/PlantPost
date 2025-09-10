#!/bin/bash

export TWD=$(pwd)
export TEST_TWD=$(pwd)
set -a
. ./config/pgtest.env
. ./config/miniotest.env

go test -cover -coverprofile=out/coverage.out -p=2 $(go list ./... | grep -v ./internal/view | grep -v ./cmd | grep -v ./internal/api | grep -v ./internal/testutils) -tags=integration
