GO ?= go
GO_BUILDFLAGS = -v
GO_TESTFLAGS = -cover

GO_BUILDDIR = ./build
GO_SRCDIRS = cli lib testutils
GO_PACKAGE_PREFIX = github.com/DirtyHairy/csync
GO_PACKAGES = \
	lib/storage \
	lib/storage/types \
	lib/storage/factory \
	lib/storage/config \
	lib/storage/local \
	lib/storage/local/config \
	lib/repository \
	lib/repository/registry \
	lib/repository/local_config \
	lib/repository/registry/config \
	lib/sync \
	lib/environment \
	lib/environment/config \
	lib/cmd/push \
	cli \
	cli/csync \
	testutils
GO_DEPENDENCIES = github.com/blang/semver

GO_DEBUG_MAIN = github.com/DirtyHairy/csync/cli/csync
GO_DEBUG_BINARY = ./csync_debug

GIT = git
GIT_COMMITFLAGS = -a

GARBAGE = $(GO_BUILDDIR) $(GO_DEBUG_BINARY)

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

godebug:
	GOPATH="`pwd`/$(GO_BUILDDIR):$$GOPATH" \
		godebug build \
		-instrument `for i in $(packages); do echo -n $$i,; done` \
		-o $(GO_DEBUG_BINARY) $(GO_DEBUG_MAIN)

godebug_run: godebug
	$(GO_DEBUG_BINARY)

godebug_test:
	@if test -z "$(PKG)"; then echo you need to set PKG to the package to test; exit 1; fi
	GOPATH="`pwd`/$(GO_BUILDDIR):$$GOPATH" \
		godebug test \
		-instrument `for i in $(packages); do echo -n $$i,; done` \
		$(GO_PACKAGE_PREFIX)/$(PKG)

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
