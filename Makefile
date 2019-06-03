GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0 

.PHONY : fmt
fmt :
	@go fmt \
		./...

.PHONY : test
test :
	@go test \
		-v \
		-race \
		./...
.PHONY : clean
clean :
	@rm -rf target/

.PHONY : build
build: clean
build: export CGO_ENABLED=1
build :
	@mkdir -p target
	@go build \
		-o target/server \
		cmd/server/main.go
