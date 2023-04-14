.PHONY: clean
clean:
	rm -rf ./dist cover.out

.PHONY: test
test:
	go test \
		-race \
		-covermode atomic \
		-coverprofile=cover.out \
		./ ./provider ./updater
	go tool cover -func cover.out

.PHONY: code-lint
code-lint:
	golangci-lint run

.PHONY: integration-test
integration-test:
	go run test/check_updates/gitlab/main.go
	go run test/check_updates/github/main.go
	go run test/update/gitlab/main.go
	go run test/update/github/main.go
