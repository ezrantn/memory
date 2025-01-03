test:
	@go test -v .

test-cov:
	@go test -coverprofile=memory.out

cov: test-cov
	@go tool cover -html=memory.out

fmt:
	@go fmt .