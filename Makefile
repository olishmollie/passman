build: main.go passman/
	go build -o bin/passman -ldflags "-X main.rootName=.passman -X main.version=$(VERSION)"

build_test: main.go passman/
	go build -o test/passman -ldflags "-X main.rootName=.passman_test"