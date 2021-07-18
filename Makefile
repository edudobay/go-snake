all: go-snake

.PHONY: go-snake init fmt

go-snake:
	go build -o go-snake

init:
	git config core.hooksPath .githooks

fmt:
	./fmt fix
