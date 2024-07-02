BINARIES=\
typeshift \


all: $(BINARIES)

typeshift: $@.go
	go build $@.go

test:
	go test

clean:
	rm -f $(BINARIES)
