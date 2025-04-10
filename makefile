COVERDIR:=out
COVERAGE_FILE:=$(COVERDIR)/coverage.out

SCRIPTS:=./scripts
INTEGRATION_TESTS:=$(SCRIPTS)/integration_tests.sh

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test -cover -coverprofile=$(COVERAGE_FILE) ./...

.PHONY: test-integration
test-integration:
	$(INTEGRATION_TESTS)

.PHONY: show-coverage
show-coverage:
	go tool cover -html=$(COVERAGE_FILE)
