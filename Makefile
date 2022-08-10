##
# urlredirector
#
# @file
# @version 0.1

VER = latest
TAG := sigmonsays/urlredirector:$(VER)


.PHONY: install

TOPDIR = $(shell pwd)

export GOWORKSPACE := $(shell pwd)
export GOBIN := $(GOWORKSPACE)/bin
export GO111MODULE := on

GO_BINS =
GO_BINS += urlredirectord

all:
	$(MAKE) compile

help:
	#
	# docker        build docker image
	# dockerpush    push docker image to $(TAG)
	#

docker:
	docker build -t $(TAG) .
dockerpush:
	docker push $(TAG)



compile:
	mkdir -p tmp
	mkdir -p $(GOBIN)
	go build -o $(GOBIN)/urlredirectord ./urlredirectord

install:
	mkdir -p $(DESTDIR)/$(INSTALL_PREFIX)/bin/
	$(MAKE) install-bins

install-bins: $(addprefix installbin-, $(GO_BINS))

$(addprefix installbin-, $(GO_BINS)):
	$(eval BIN=$(@:installbin-%=%))
	cp -v $(GOBIN)/$(BIN) $(DESTDIR)/$(INSTALL_PREFIX)/bin/

# end
