BINARIES=\
bogusgroup \


all: $(BINARIES)

bogusgroup: $@.go
	go build $@.go

test:
	go test

clean:
	rm -f $(BINARIES)
