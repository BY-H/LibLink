.PHONY: build
build:
	@go build -o _output/liblink cmd/liblink/main.go
	@mkdir -p _output/static
	@cp -r static/* _output/static/

.PHONY: format
format:
	@gofmt -s -w ./

.PHONY: clean
clean:
	@rm -vrf _output
