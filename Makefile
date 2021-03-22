# Default GO_BIN to Go binary in PATH
GO_BIN				?= go
DOCKER_BIN			?= docker

TEST_PATTERN ?=.
TEST_OPTIONS ?=
SOURCE_FILES ?= ./...

TEST_FLAGS += -failfast
TEST_FLAGS += -race

GO_TEST ?= test $(TEST_OPTIONS) $(TEST_FLAGS) $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=10m

.PHONY: go-get
go-get:
	@printf '\n================================================================\n'
	@printf 'Target: go-get'
	@printf '\n================================================================\n'
	$(GO_BIN) mod vendor
	@echo '[go-get] Done.'

.PHONY: test-coverage
test-coverage: TEST_FLAGS += -covermode=atomic -coverprofile=coverage.out
test-coverage: go-get
	@printf '\n================================================================\n'
	@printf 'Target: test-coverage'
	@printf '\n================================================================\n'
	@echo '[test] Testing packages: $(SOURCE_FILES)'
	$(GO_BIN) $(GO_TEST)

.PHONY: docker/go-get
docker/go-get:
	@($(DOCKER_BIN) run --rm -it -v ${CURDIR}:/app -w /app golang:1.14 make go-get)

.PHONY: quick-start
quick-start: docker/go-get
	@printf '\n================================================================\n'
	@printf 'Target: quick-start'
	@printf '\n================================================================\n'
	$(DOCKER_BIN) run --rm -it -v ${CURDIR}:/app -w /app golang:1.14 make start

bin/rover: go-get
	$(GO_BIN) build -o ${CURDIR}/bin/server ./cmd/server

build: bin/rover

.PHONY: test
test: go-get
	@printf '\n================================================================\n'
	@printf 'Target: test'
	@printf '\n================================================================\n'
	$(GO_BIN) $(GO_TEST)

.PHONY: docker/test
docker/test:
	@($(DOCKER_BIN) run --rm -it -v ${CURDIR}:/app -w /app golang:1.14 make test)

.PHONY: start
start:
	$(GO_BIN) run ./cmd/server/

