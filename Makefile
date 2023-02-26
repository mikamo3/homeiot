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

build: clean
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINDIR)/$(PROJECT_NAME)_scan cmd/scan/main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINDIR)/$(PROJECT_NAME)_store cmd/store/main.go

build-rpi: clean
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BINDIR)/scan cmd/scan/main.go
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(BINDIR)/store cmd/store/main.go
