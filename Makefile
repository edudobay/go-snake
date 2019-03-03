all: go-snake

.PHONY: go-snake fmt

go-snake:
	source goenv/activate && \
	    go build -o go-snake

fmt:
	find . \! -path './goenv/*' -name '*.go' -print0 | \
	    xargs -0 gofmt -l -w
