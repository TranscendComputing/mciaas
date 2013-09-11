all: deps
	@./scripts/build.sh

deps:
	@echo "--> Installing dependencies"
	@go get -d -v ./...
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs -n1 go get -d

clean:
	@rm -rf mciaas bin/ local/ pkg/ src/

format:
	go fmt ./...

test: deps
	@echo "--> Testing MCIaaS..."
	go test ./...

.PHONY: all deps format test
