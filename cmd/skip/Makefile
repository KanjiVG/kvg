BINARIES=\
skip skip.json \


all: $(BINARIES)

skip: $@.go
	go build $@.go

skip.json: make-skip-json.pl
	perl make-skip-json.pl

test:
	go test

clean:
	rm -f $(BINARIES)
