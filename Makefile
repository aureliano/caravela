.PHONY: clean
clean:
	rm -rf ./dist cover.out

.PHONY: test
test:
	go test \
        -v \
        --count=1 \
        -coverprofile cover.out \
		./http
	go tool cover -func cover.out

.PHONY: code-lint
code-lint:
	go vet -json ./
