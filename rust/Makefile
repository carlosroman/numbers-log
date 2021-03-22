CURRENT_USERID 		?= $(shell id -u)
CURRENT_GROUP		?= $(shell id -g)

CARGO_BIN		?= cargo
DOCKER_BIN		?= docker
RUST_VERSION		?= 1.50.0

DOCKER_TTY				?= -t
RUST_DOCKER_BUILD_ARG	+= --build-arg CURRENT_GROUP=$(CURRENT_GROUP)
RUST_DOCKER_BUILD_ARG	+= --build-arg CURRENT_USERID=$(CURRENT_USERID)
RUST_DOCKER_BUILD_ARG	+= --build-arg RUST_VERSION=$(RUST_VERSION)
RUST_DOCKER_CMD			?= $(DOCKER_BIN) run --rm --user "$(CURRENT_USERID)":"$(CURRENT_GROUP)" $(DOCKER_TTY) -v $(CURDIR):/app ci/rust
DOCKER_TARGET_CMD		= $(word 2, $(subst /, ,$(@)))
CI_CMD					?= $(RUST_DOCKER_CMD) make $(DOCKER_TARGET_CMD)
CI_DEP_BUILD_CMD		?= $(DOCKER_BIN) build $(RUST_DOCKER_BUILD_ARG) -t $(CI_DOCKER_TAG) -f $(CURDIR)/build/ci/Dockerfile .
CI_DOCKER_TAG			?= ci/rust

.PHONY : fmt
fmt :
	@(rustfmt src/*)

.PHONY : test
test :
	@($(CARGO_BIN) test)

.PHONY : start
start :
	@($(CARGO_BIN) run)

.PHONY : quick-start
quick-start : DOCKER_TTY += -i
quick-start : ci/dep
	@printf '\n================================================================\n'
	@printf 'Target: quick-start'
	@printf '\n================================================================\n'
	@($(RUST_DOCKER_CMD) make start)

target/release/rover :
	$(CARGO_BIN) build --release

build: target/release/rover

.PHONY : lint
lint :
	@($(CARGO_BIN) clippy)

.PHONY : all
all : fmt lint test build

.PHONY : ci/test
ci/test : ci/dep
	@printf 'CMD: $(RUST_DOCKER_CMD)\n'
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
	@($(CI_DEP_BUILD_CMD))
