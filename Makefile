NAME=$(lastword $(subst /, ,$(abspath .)))
VERSION=$(shell git.exe describe --tags 2>nul || echo v0.0.0)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"

ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=SET
else
    SET=export
endif

all:
	go fmt $(foreach X,$(wildcard internal/*),&& cd $(X) && go fmt && cd ../..)
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

_package:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(EXT)

package:
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=386"   && $(MAKE) _package EXT=
	$(SET) "GOOS=linux"   && $(SET) "GOARCH=amd64" && $(MAKE) _package EXT=
	$(SET) "GOOS=windows" && $(SET) "GOARCH=386"   && $(MAKE) _package EXT=.exe
	$(SET) "GOOS=windows" && $(SET) "GOARCH=amd64" && $(MAKE) _package EXT=.exe

manifest:
	make-scoop-manifest *-windows-*.zip > $(NAME).json

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
