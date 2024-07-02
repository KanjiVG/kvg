BINARIES=\
empty-path \


all: $(BINARIES)

empty-path: $@.go
	go build $@.go

test:
	go test

clean:
	rm -f $(BINARIES)
