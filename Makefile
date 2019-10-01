build:
	GOOS=darwin GOARCH=386   go build -o ./bin/mac/kns ./kns.go
	GOOS=linux  GOARCH=amd64 go build -o ./bin/amd64/kns ./kns.go

install_macos:
	cp -p ./bin/mac/kns /usr/local/bin/

install_linux:
	cp -p ./bin/amd64/kns /usr/local/bin/