# Default GO_BIN to Go binary in PATH
GO_BIN				?= go
GO_VERSION			?= 1.17
DOCKER_BIN			?= docker

TEST_PATTERN ?=.
TEST_OPTIONS ?=
SOURCE_FILES ?= ./...

TEST_FLAGS += -failfast
TEST_FLAGS += -race

GO_TEST ?= test $(TEST_OPTIONS) $(TEST_FLAGS) $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=10m

DOCKER_TTY				?= -t
DOCKER_CMD				?= $(DOCKER_BIN) run --rm $(DOCKER_TTY) -v ${CURDIR}:/app -w /app golang:$(GO_VERSION)
CI_CMD					?= $(DOCKER_CMD) make $(DOCKER_TARGET_CMD)
DOCKER_TARGET_CMD		= $(word 2, $(subst /, ,$(@)))

.PHONY: go-get
go-get:
	@printf '\n================================================================\n'
	@printf 'Target: go-get'
	@printf '\n================================================================\n'
	$(GO_BIN) mod vendor
	@echo '[go-get] Done.'

.PHONY: dep
dep: go-get

.PHONY: test-coverage
test-coverage: TEST_FLAGS += -covermode=atomic -coverprofile=coverage.out
test-coverage: dep
	@printf '\n================================================================\n'
	@printf 'Target: test-coverage'
	@printf '\n================================================================\n'
	@echo '[test] Testing packages: $(SOURCE_FILES)'
	$(GO_BIN) $(GO_TEST)

.PHONY: quick-start
quick-start: ci/dep
	@printf '\n================================================================\n'
	@printf 'Target: quick-start'
	@printf '\n================================================================\n'
	$(DOCKER_CMD) make start

bin/server: dep
	$(GO_BIN) build -o ${CURDIR}/bin/server ./cmd/server

build: bin/server

.PHONY: test
test: dep
	@printf '\n================================================================\n'
	@printf 'Target: test'
	@printf '\n================================================================\n'
	$(GO_BIN) $(GO_TEST)

.PHONY: start
start:
	$(GO_BIN) run ./cmd/server/

.PHONY : ci/test
ci/test : ci/dep
	@($(CI_CMD))

.PHONY : ci/test-coverage
ci/test-coverage : ci/dep
	@($(CI_CMD))

.PHONY : ci/build
ci/build : ci/dep
	@($(CI_CMD))

.PHONY : ci/lint
ci/lint : ci/dep
	@($(CI_CMD))

.PHONY : ci/all
ci/all : ci/dep
	@($(CI_CMD))

.PHONY : ci/dep
ci/dep :
	@(echo "noop")
