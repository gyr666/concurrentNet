BUILD=go build
build:
	$(BUILD) $(FILE)
free:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BUILD)
linux:
	GOOS=linux GOARCH=amd64 $(BUILD)
run:
	./main
install:
	mv concurrentNet /bin
clean:
	rm main
