build:
	go build -o ./bin/kns ./kns.go
install:
	go build -o ./bin/kns ./kns.go && cp -p ./bin/kns /usr/local/bin/
