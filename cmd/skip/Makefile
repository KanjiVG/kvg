BINARIES=\
skip \


all: $(BINARIES) skip.json 

skip: $@.go
	go build $@.go

skip.json: make-skip-json.pl
	perl make-skip-json.pl

test:
	go test

clean:
	rm -f $(BINARIES)
