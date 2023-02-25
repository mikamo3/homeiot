build-rpi:
	GOOS=linux GOARCH=arm go build -o scan cmd/scan/main.go
	GOOS=linux GOARCH=arm go build -o store cmd/store/main.go