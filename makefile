COVERDIR:=out
COVERAGE_FILE:=$(COVERDIR)/coverage.out

SCRIPTS:=./scripts
INTEGRATION_TESTS:=$(SCRIPTS)/integration_tests.sh

COMPOSEFILE:=./deployments/docker-compose.yaml
COMPOSEFILE_DEV:=./deployments/docker-compose.dev.yaml

PARSEDEPTH:= 10
SWAGGER:=./cmd/docs
API_APP:=./cmd/api/main.go
API_DIR:=./cmd/api
API_BUILD := ./api.bin
TEMPL_DIR:=./internal/view


.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test $$(go list ./... | grep -v ./internal/view | grep -v ./cmd) -cover -coverprofile=$(COVERAGE_FILE) 

.PHONY: test-integration
test-integration:
	$(INTEGRATION_TESTS)

.PHONY: show-coverage
show-coverage:
	go tool cover -html=$(COVERAGE_FILE)

.PHONY: api-build
api-build:
	swag init --parseInternal --parseDependency --parseDepth $(PARSEDEPTH) -g $(API_APP) -o $(SWAGGER)
	tailwindcss -o ./internal/view/static/css/tailwind.css --minify
	go tool templ generate -path $(TEMPL_DIR)
	npm install && npm run build
	go build -o $(API_BUILD) $(API_DIR)

.PHONY: api-run
api-run:
	$(API_BUILD)

.PHONY: api
api: api-build api-run

.PHONY: dev-up
dev-up:
	docker compose -f $(COMPOSEFILE_DEV) up 

.PHONY: dev-upd
dev-upd:
	docker compose -f $(COMPOSEFILE_DEV) up -d 

.PHONY: dev-update
dev-update:
	docker compose -f $(COMPOSEFILE_DEV) up --build

.PHONY: dev-down
dev-down:
	docker compose -f $(COMPOSEFILE_DEV) down 

.PHONY: up
up:
	docker compose -f $(COMPOSEFILE) up 

.PHONY: upd
upd:
	docker compose -f $(COMPOSEFILE) up -d 

.PHONY: down
down:
	docker compose -f $(COMPOSEFILE) down 
