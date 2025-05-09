run: build
	go run *

build: format
	go build

format:
	gofmt -w .

clean:
	rm boid
