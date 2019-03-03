all: go-snake

.PHONY: go-snake init fmt

go-snake:
	source goenv/activate && \
	    go build -o go-snake

init:
	git config core.hooksPath .githooks

fmt:
	./fmt fix
