.PHONY: fmt lint vulncheck

fmt:
	go tool gofumpt -w .

lint:
	go tool staticcheck -checks=all -show-ignored -tests ./...

vulncheck:
	go tool govulncheck ./...
