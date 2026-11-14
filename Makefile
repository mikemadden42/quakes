build:
	goreleaser build --snapshot --clean

check:
	golangci-lint run

snapshot:
	goreleaser release --snapshot --clean

.PHONY: build check snapshot
