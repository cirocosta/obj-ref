run: build
	./main -o yaml pods foo

build:
	go build -v -i -o main
