PREFIX := /usr/local
PKGNAME := chhoto

build:
	go build -ldflags="-s -w" -o ${PKGNAME}

install: build
	install -Dm755 $(PKGNAME) "$(DESTDIR)$(PREFIX)/bin/$(PKGNAME)"
	install -Dm644 $(PKGNAME).1 "$(DESTDIR)$(PREFIX)/man/man1/$(PKGNAME).1"

uninstall:
	rm -f "$(DESTDIR)$(PREFIX)/bin/$(PKGNAME)"
	rm -f "$(DESTDIR)$(PREFIX)/man/man1/$(PKGNAME).1"

conf_tag := $(shell cat internal/config.go | sed -rn 's/^const version = "(.+)"$$/\1/p')
last_tag := $(shell git tag -l | tail -1)
bumped := $(shell git log -1 --pretty=%B | grep "build: Bumped version to " | wc -l)
tag:
ifneq (${conf_tag}, ${last_tag})
ifeq (${bumped}, 1)
	git tag ${conf_tag} -m "Version ${conf_tag}"
endif
endif

aur: build tag
	git push
	tar --transform 's/.*\///g' -czf $(PKGNAME).tar.gz $(PKGNAME) $(PKGNAME).1

.PHONY: build install uninstall aur tag
