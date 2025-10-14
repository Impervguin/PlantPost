#!/bin/bash

COVERDIR=out
COVERAGE_FILE=$COVERDIR/coverage.out

go test $(go list ./... | grep -v ./internal/ | grep -v ./cmd) -cover -coverprofile=$COVERAGE_FILE -tags=unit