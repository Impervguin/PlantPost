#!/bin/bash

export TWD=$(pwd)
export TEST_TWD=$(pwd)
set -a
. ./config/pgtest.env
. ./config/miniotest.env

go test -cover -coverprofile=out/coverage.out ./... -tags=integration
