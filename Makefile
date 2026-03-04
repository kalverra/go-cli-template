.PHONY: goreleaser lint test test_race bench

goreleaser:
	goreleaser release --clean --snapshot

lint:
	go fix ./...
	golangci-lint run ./... --fix

test:
	go tool gotestsum -- -cover ./...

test_race:
	go tool gotestsum -- -race ./...

bench:
	go test -bench=. -benchmem -run=^$$ ./... -cpu=2,4,8
