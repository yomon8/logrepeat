BIN      := logrepeat
VERSION  := $(shell git describe --tags)

all: build

test: deps build
	go test ./...

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get github.com/mitchellh/gox
	go get github.com/fatih/color

lint: deps
	go vet ./...
	golint -set_exit_status ./...

crossbuild:
	rm -fR ./pkg && mkdir ./pkg ;\
		gox \
		-output "./pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" \
		-ldflags "-X main.version=$(VERSION)" \
		./cmd/...

build:
	go build -o $(BIN) -ldflags "-X main.version=$(VERSION)" ./cmd/...

linuxbuild:
	GOOS=linux GOARCH=amd64 make build

clean:
	go clean
