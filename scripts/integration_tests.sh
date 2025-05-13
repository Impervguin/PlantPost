#!/bin/bash

export TWD=$(pwd)
export TEST_TWD=$(pwd)
set -a
. ./config/pgtest.env
. ./config/miniotest.env

go test -cover -coverprofile=out/coverage.out $(go list ./... | grep -v ./internal/view | grep -v ./cmd) -tags=integration
