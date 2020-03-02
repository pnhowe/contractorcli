all: contractorcli

contractorcli: main.go contractor/* cmd/*
	go build -ldflags "-linkmode external -extldflags -static" -o contractorcli -a main.go 
