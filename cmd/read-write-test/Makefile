BINARIES=\
read-write-test \


all: $(BINARIES)

read-write-test: $@.go
	go build $@.go

test:
	go test

clean:
	rm -f $(BINARIES)
