GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0 

.PHONY : fmt
fmt :
	@go fmt \
		./...

.PHONY : test-ci
test-ci: clean
	@mkdir -p target
	@chmod +x scripts/coverage.sh
	@scripts/coverage.sh

.PHONY : test
test :
	@go test \
		-v \
		-race \
		./...
.PHONY : clean
clean :
	@rm -rf target/

.PHONY : lint-install
lint-install :
	@go get \
	    github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0

.PHONY : lint
lint :
	@golangci-lint run

.PHONY : build
build: clean
build: export CGO_ENABLED=1
build :
	@mkdir -p target
	@go build \
		-o target/server \
		cmd/server/main.go
