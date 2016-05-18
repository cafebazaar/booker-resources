GIT ?= git
GO ?= go
DOWNLOADER ?= curl -o
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags --exact-match --abbrev=0 --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
TARGET := github.com/cafebazaar/booker-reservation
TARGET_DOTS := "$(shell echo $(TARGET) | sed -e "s/\//./g")"
LD_FLAGS := -X $(TARGET)/common.Version=$(VERSION) -X $(TARGET)/common.BuildTime=$(BUILD_TIME)
FORMAT := '{{ join .Deps " " }}'

.PHONY: help clean dependencies docker
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  reservation  to build the main binary for current platform"
	@echo "  dependencies to go install the dependencies"
	@echo "  docker       to build the docker image"
	@echo "  clean        to remove generated files"

clean:
	rm -f reservation
	rm -f proto/internal.pb.gw.go proto/internal.pb.go proto/public.pb.go proto/common.pb.go
	rm -f proto/common.proto

proto/common.proto:
	cd proto; $(DOWNLOADER) common.proto "https://raw.githubusercontent.com/cafebazaar/booker-apiserver/master/proto/common.proto"

proto/internal.pb.gw.go proto/internal.pb.go proto/public.pb.go proto/common.pb.go: proto/internal.proto proto/public.proto proto/common.proto
	cd proto; go generate -v .

dependencies:
	$(GO) list -f=$(FORMAT) $(TARGET) | xargs $(GO) install

reservation: proto/internal.pb.gw.go proto/internal.pb.go proto/public.pb.go proto/common.pb.go main.go */*.go
	$(GO) build -o="reservation" -ldflags="$(LD_FLAGS)" $(TARGET)

docker: reservation Dockerfile
	docker build -t $(DOCKER_IMAGE) .
