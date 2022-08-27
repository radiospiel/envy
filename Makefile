.PHONY: doc

# we currently don't need MACHTYPE
# MACHTYPE:=$(shell set | grep ^MACHTYPE= | sed s-.*=--)

default: build

test: build test

all: build test build-all

# -- go specific --------------------------------------------------------------

.PHONY: bin/envy.go.bin
build: bin/envy.go.bin

bin/envy.go.bin: Makefile
	cd src/golang && go build -o ../../$@

test:
	ENVY=bin/envy.go spec/bin/roundup spec/*-test.sh

# --- releases ----------------------------------------------------------------

# builds and packs binaries for all supported platforms
#
# $MACHTYPE                | $GOOS  | $GOARCH
#
# x86_64-apple-darwin18    | darwin | amd64
# x86_64-pc-linux-gnu      | linux  | amd64
# arm64-apple-darwin21     | darwin | amd64
#
build-all: .dependencies
	cd src/golang && GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../../bin/envy.x86_64-apple-darwin.bin
	cd src/golang && GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../../bin/envy.arm64-apple-darwin.bin
	cd src/golang && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../../bin/envy.x86_64-pc-linux-gnu.bin
	upx bin/envy.*arm64*.bin
	upx bin/envy.*x86*.bin

prerelease: build-all
	scripts/prerelease

release: build-all
	scripts/release

# --- local installation ------------------------------------------------------

install: bin/envy.go.bin
	install bin/envy.go.bin /usr/local/bin/envy
	mkdir -p /usr/local/share/man/man1/
	[ -f doc/envy.1 ] && install doc/*.1 /usr/local/share/man/man1/ || true

doc: doc/envy.1

doc/envy.1: README.md
	@mkdir -p doc
	@which -s ronn || (echo "Please install ronn: gem install ronn" && false)
	ronn --pipe --roff README.md > doc/envy.1
	ronn --pipe --html README.md > doc/envy.1.html
