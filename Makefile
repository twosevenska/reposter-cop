BUILD_CONTEXT := ./build

.PHONY: go-test
go-test:
	@echo "Run all project tests..."
	go test -v -p 1 ./...

.PHONY: go-build
go-build: go-get
	@echo "Build project binaries..."
	CGO_ENABLED=0 go build -v -o $(BUILD_CONTEXT)/reposter-cop cop.go

.PHONY: go-get
go-get:
	@echo "Fetch project dependencies..."
	GO111MODULE=on go get -u ./...
	GO111MODULE=on go mod vendor

.PHONY: clean
clean:
	rm -rf $(BUILD_CONTEXT)/resposter-cop

.PHONY: setup
setup: clean go-get go-build

.PHONY: test
test: go-get go-test

.PHONY: run
run:
	GO111MODULE=on go run cop.go
