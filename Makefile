PROG= check-receiver

all: bin

bin:
	go build

deb: bin
	debuild -uc -us -i

strip:
	strip --strip-all $(PROG)

test:
	go test -gocheck.v

format:
	gofmt -s -w *.go

clean:
	rm -f $(PROG)
