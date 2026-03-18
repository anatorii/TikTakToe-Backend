.PHONY: server clean

server:
	go build -o server ./cmd/tiktaktoe/main.go
	./server

clean:
	rm -f ./server

