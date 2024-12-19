MODULE := $(shell head -n 1 go.mod | sed -E 's/^module //')

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
