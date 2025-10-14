#!/bin/bash

export TWD=$(pwd)
export TEST_TWD=$(pwd)
export TEST_POSTGRES_EXTERNAL=true
export TEST_MINIO_EXTERNAL=true
set -a
. ./config/pgtest.env
. ./config/miniotest.env

go test -cover -coverprofile=out/coverage.out -p=1 $(go list ./... | grep -v ./internal/view | grep -v ./e2e | grep -v ./cmd | grep -v ./internal/api | grep -v ./internal/testutils) -tags=integration
