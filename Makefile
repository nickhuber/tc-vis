.PHONY: all install clean

all: tc-vis

install: tc-vis
	install tc-vis /usr/local/bin/

clean:
	go clean

tc-vis: tc-vis.go
	go build
