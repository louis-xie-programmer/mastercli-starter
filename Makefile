APP=mastercli

.PHONY: build run lint test clean

build:
	go build -o $(APP) ./cmd/mastercli

run:
	go run ./cmd/mastercli --help

clean:
	rm -f $(APP)

