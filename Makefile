ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    WHICH=where.exe
    DEL=del
    NUL=nul
else
    SET=export
    WHICH=which
    DEL=rm
    NUL=/dev/null
endif

ifndef GO
    SUPPORTGO=go1.20.14
    GO:=$(shell $(WHICH) $(SUPPORTGO) 2>$(NUL) || echo go)
endif

NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT:=-ldflags "-s -w -X main.version=$(VERSION)"
EXE:=$(shell $(GO) env GOEXE)

all:
	$(GO) fmt ./...
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)

_dist:
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXE)

dist:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _dist
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _dist

release:
	$(GO) run github.com/hymkor/latest-notes@latest | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

manifest:
	$(GO) run github.com/hymkor/make-scoop-manifest@master -all *-windows-*.zip > $(NAME).json

test:
	$(SET) "GOFILE=const.go" && $(SET) "GOLINE=3" && $(SET) "GOPACKAGE=dos" && cd example && "../go-importconst" -d

test-lowercamel:
	$(SET) "GOFILE=const.go" && $(SET) "GOLINE=3" && $(SET) "GOPACKAGE=dos" && cd example && "../go-importconst" -d -lowercamel

test-uppercamel:
	$(SET) "GOFILE=const.go" && $(SET) "GOLINE=3" && $(SET) "GOPACKAGE=dos" && cd example && "../go-importconst" -d -uppercamel

test-underscore:
	$(SET) "GOFILE=const.go" && $(SET) "GOLINE=3" && $(SET) "GOPACKAGE=dos" && cd example && "../go-importconst" -d -prefix _

test2:
	$(SET) "GOFILE=const.go" && $(SET) "GOLINE=3" && $(SET) "GOPACKAGE=dos" && cd example2 && "../go-importconst" -d
