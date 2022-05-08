# Cross-compile for Windows
export GOOS := windows
export GOARCH := amd64

# Stop polluting my $HOME
export GOMODCACHE := $(realpath cache/pkg/mod)

cmdgo.exe: cmdgo.go
	go build -o $@ $<

cmdgo.go: deps

deps:
	go get golang.org/x/sys/windows

fmt:
	go fmt

clean:
	-rm -f cmdgo.exe
	-chmod -R +w $(GOMODCACHE); rm -rf $(GOMODCACHE)/{cache,golang.org}

.PHONY: deps fmt clean
