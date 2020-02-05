BULID=go build
Default:
	go build main.go
Free:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BUILD) main.go
Lin:
	GOOS=linux GOARCH=amd64 $(BUILD) main.go
Run:
	./main
Clean:
	rm main
