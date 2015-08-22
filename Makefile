GO ?= go
GO_BUILDFLAGS = -v
GO_TESTFLAGS = -cover

GO_BUILDDIR = ./build
GO_SRCDIRS = cmd lib
GO_PACKAGE_PREFIX = github.com/DirtyHairy/csync
GO_PACKAGES = cmd/csync lib/storage lib/storage/local
GO_DEPENDENCIES =

GIT = git
GIT_COMMITFLAGS = -a

GARBAGE = $(GO_BUILDDIR)

packages = $(GO_PACKAGES:%=$(GO_PACKAGE_PREFIX)/$(GO_SRCDIR)/%)
execute_go = GOPATH=`pwd`/$(GO_BUILDDIR) $(GO) $(1) $(2) $(packages)

all: install

install: $(GO_BUILDDIR)
	$(call execute_go,install,$(GO_BUILDFLAGS))

fmt: $(GO_BUILDDIR)
	$(call execute_go,fmt)

goclean: $(GO_BUILDDIR)
	$(call execute_go,clean)

test: $(GO_BUILDDIR)
	$(call execute_go,test,$(GO_TESTFLAGS))

vet: $(GO_BUILDDIR)
	$(call execute_go,vet)

commit: fmt
	$(GIT) commit $(GIT_COMMITFLAGS)

$(GO_BUILDDIR):
	mkdir -p ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)
	for srcdir in $(GO_SRCDIRS); \
	    do \
	    	ln -s `pwd`/$$srcdir ./$(GO_BUILDDIR)/src/$(GO_PACKAGE_PREFIX)/$$srcdir; \
	    done
	if test -n "$(GO_DEPENDENCIES)"; then GOPATH=`pwd`/$(GO_BUILDDIR) $(GO) get $(GO_DEPENDENCIES); fi

clean:
	-rm -fr $(GARBAGE)

.PHONY: clean all install fmt goclean test
