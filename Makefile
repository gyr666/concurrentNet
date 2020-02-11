BUILD=go build
FILE=main.go
Default:
	$(BUILD) $(FILE)
Free:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BUILD) $(FILE)
Lin:
	GOOS=linux GOARCH=amd64 $(BUILD) $(FILE)
Run:
	./main
Clean:
	rm main
