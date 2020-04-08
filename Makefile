.POSIX:

all: sreader

sreader: sreader.go
	go build $^

run: sreader.go
	go run $^
