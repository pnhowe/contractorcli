VERSION = $(shell git tag -l | tail -n1)
GIT_VERSION = $(shell git rev-list -1 HEAD)

all: contractorcli

install: contractorcli
	mkdir -p $(DESTDIR)/usr/bin
	install -m 755 contractorcli $(DESTDIR)/usr/bin

version:
	echo $(VERSION)

contractorcli: main.go cmd/* go.mod go.sum
	go build -ldflags "-linkmode external -extldflags -static -X 'github.com/t3kton/contractorcli/cmd.version=${VERSION}' -X 'github.com/t3kton/contractorcli/cmd.gitVersion=${GIT_VERSION}'" -o contractorcli -a main.go

clean:
	$(RM) contractorcli
	$(RM) dpkg
	dh_clean || true

dist-clean: clean

.PHONY:: version clean dist-clean

dpkg-blueprints:
	echo ubuntu-focal-base

dpkg-requires:
	echo dpkg-dev debhelper golang-1.13 golang

dpkg:
	dpkg-buildpackage -b -us -uc
	touch dpkg

dpkg-file:
	echo $(shell ls ../contractorcli_*.deb):focal

.PHONY:: dpkg-blueprints dpkg-requires dpkg dpkg-file
