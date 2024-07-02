BINARIES=\
missing-stroke \


all: $(BINARIES)

missing-stroke: $@.go
	go build $@.go

test:
	go test

clean:
	rm -f $(BINARIES)
