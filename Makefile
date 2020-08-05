check: test lint vet fmt-check
ci-check: test vet fmt-check

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt-check:
	goimports -l *.go **/*.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

fmt:
	gofmt -w -s *.go **/*.go
	goimports -w *.go **/*.go

