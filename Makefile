# TODO
PROJECT_NAME:=homeiot
BINDIR:=bin
LDFLAGS:=-ldflags "-s"
.PHONY: clean
clean:
	@rm -rf $(BINDIR)/*

.PHONY: test
test:
	cd lib && go test

build-scan:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINDIR)/$(PROJECT_NAME)_scan cmd/scan/main.go

build-store:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINDIR)/$(PROJECT_NAME)_store cmd/store/main.go

build: clean build-scan build-store

build-rpi-scan:
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BINDIR)/$(PROJECT_NAME)_scan cmd/scan/main.go

build-rpi-store:
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BINDIR)/$(PROJECT_NAME)_store cmd/store/main.go

build-rpi: clean build-rpi-scan build-rpi-store
