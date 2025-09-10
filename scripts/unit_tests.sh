#!/bin/bash

export ALLURE_OUTPUT_PATH=$PWD
export ALLURE_OUTPUT_DIR=allure-results

COVERDIR=out
COVERAGE_FILE=$COVERDIR/coverage.out

go test $(go list ./... | grep -v ./internal/view | grep -v ./cmd | grep -v ./internal/api | grep -v ./internal/testutils) -cover -coverprofile=$COVERAGE_FILE -tags=unit