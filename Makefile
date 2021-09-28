
build:
	go build -o bin/kubectl-check

install: build
	mv bin/kubectl-check /usr/local/bin/kubectl-check