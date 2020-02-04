all: install

install:
	go install goplugin

clean:
	go clean ./...