#!/bin/bash

COVERDIR=out
COVERAGE_FILE=$COVERDIR/coverage.out

go test $(go list ./... | grep -v ./internal/view | grep -v ./cmd | grep -v ./e2e | grep -v ./internal/api | grep -v ./internal/testutils) -cover -coverprofile=$COVERAGE_FILE -tags=unit