PROG= nagios-receiver

all: bin

bin:
	go build

strip:
	strip --strip-all $(PROG)

format:
	gofmt -s -w *.go

clean:
	rm -f $(PROG)
