.PHONY: mod test run post

all: test

mod:
	go mod tidy
	go mod vendor

test: mod
	go test -v -count=1 -race ./...

run: mod
	go run -v handler.go main.go

post:
	curl -d '{"id":"id-1234", "txt":"test"}' \
		-H "Content-type: application/json" \
		"http://localhost:8082/decorate"