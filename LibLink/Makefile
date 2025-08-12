.PHONY: build
build:
	@go build -o _output/liblink cmd/liblink/main.go

.PHONY: format
format:
	@gofmt -s -w ./

.PHONY: clean
clean:
	@rm -vrf _output
