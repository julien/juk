MODULE := $(shell head -n 1 go.mod | sed -E 's/^module //')

all:
	go vet ./... && go run .

fmt:
	@goimports -local $(MODULE) -w .

test:
	@go test ./... \
		-count 1 \
		-cover \
		-covermode atomic \
		-coverprofile coverage.cov \
		-v

.PHONY: fmt test
