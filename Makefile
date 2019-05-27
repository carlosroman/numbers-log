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
