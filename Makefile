.PHONY: build
build:
	CGO_ENABLED=0 cd cmd && go build -o ../bin/deilephila -ldflags "-w -s" ./deilephila

.PHONY: build-plugin
build-plugin:
	go build -buildmode=plugin -o ./bin/plugin.so ./plugin/pluginTemplate.go

.PHONY: install
install:
	GOBIN=$(PWD)/bin && go install

.PHONY: lint
lint:
	golint ./...
	cd cmd && golint ./...

.PHONY: vet
vet:
	go vet ./...
	cd cmd && go vet ./...

.PHONY: test
test: unit-test

.PHONY: unit-test
unit-test:
	go test ./...

.PHONY: integration-test
integration-test:
	go clean -testcache && cd ./test/integration && go test ./...

.PHONY: benchmark-test
benchmark-test:
	cd ./test/benchmark && go test -bench=.
