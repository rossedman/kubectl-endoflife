
.PHONY: test install

build: bin/kubectl-check
	go build -o bin/kubectl-check

install: build
	mv bin/kubectl-check /usr/local/bin/kubectl-check

test: 
	go test -coverprofile=coverage.out -race -v ./...