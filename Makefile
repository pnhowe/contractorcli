VERSION = 0.3
GIT_VERSION = $(shell git rev-list -1 HEAD)

all: contractorcli

contractorcli: main.go cmd/* go.mod go.sum
	go build -ldflags "-linkmode external -extldflags -static -X 'github.com/t3kton/contractorcli/cmd.version=${VERSION}' -X 'github.com/t3kton/contractorcli/cmd.gitVersion=${GIT_VERSION}'" -o contractorcli -a main.go

clean:
	${RM} contractorcli

.PHONY: clean
