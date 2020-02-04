all: install

install:
	go install ./plugin

clean:
	go clean ./...