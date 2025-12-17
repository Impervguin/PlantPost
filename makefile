SCRIPTS:=./scripts
INTEGRATION_TESTS:=$(SCRIPTS)/integration_tests.sh
E2E_TESTS:=$(SCRIPTS)/e2e_tests.sh

COMPOSEFILE:=./deployments/docker-compose.yaml
COMPOSEFILE_DEV:=./deployments/docker-compose.dev.yaml

PARSEDEPTH:= 10
SWAGGER:=./cmd/docs
API_APP:=./cmd/api/main.go
API_DIR:=./cmd/api
API_BUILD := ./api.bin
TEMPL_DIR:=./internal/view

ALLURE_OUTPUT_DIR:=allure-results
export ALLURE_OUTPUT_DIR
ALLURE_OUTPUT_PATH:=$(PWD)
export ALLURE_OUTPUT_PATH
ALLURE_REPORT_DIR:=$(PWD)/allure-report

.PHONY: test
test: test-unit test-integration test-e2e

.PHONY: test-unit 
test-unit: allure-clear
	$(SCRIPTS)/unit_tests.sh
	 
.PHONY: test-integration 
test-integration: allure-clear
	$(INTEGRATION_TESTS)

.PHONY: test-integration-external
test-integration-external: allure-clear
	$(SCRIPTS)/integration_external.sh

.PHONY: test-integration-parallel
test-integration-parallel: allure-clear
	$(SCRIPTS)/integration_external.sh & $(SCRIPTS)/integration_external.sh & $(SCRIPTS)/integration_external.sh & $(SCRIPTS)/integration_external.sh


.PHONY: test-e2e
test-e2e: allure-clear
	$(E2E_TESTS)

.PHONY: allure-clear
allure-clear:
	rm -rf $(ALLURE_OUTPUT_DIR)

.PHONY: allure-regen
allure-regen:
	cp -r $(ALLURE_REPORT_DIR)/history $(ALLURE_OUTPUT_DIR)
	allure generate --clean

.PHONY: allure-serve
allure-serve:
	allure open

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

update:
	docker compose -f $(COMPOSEFILE) up --build

.PHONY: down
down:
	docker compose -f $(COMPOSEFILE) down 
