
build:
	go build -o bin/kubectl-tks

install: build
	mv bin/kubectl-tks /usr/local/bin/kubectl-tks