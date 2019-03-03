all: go-snake

.PHONY: go-snake

go-snake:
	source goenv/activate && \
	    go build -o go-snake
