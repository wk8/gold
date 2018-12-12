.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -v -cover

.PHONY: cover
cover:
	rm -fv cover.* && go test -cover -coverprofile cover.out && go tool cover -html=cover.out -o cover.html && open cover.html
