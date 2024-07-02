BINARIES=\
renumber \


all: $(BINARIES)

renumber: $@.go
	go build $@.go

test:
	go test

clean:
	rm -f $(BINARIES)
