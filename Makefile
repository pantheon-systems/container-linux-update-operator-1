.PHONY:	all image clean test coverage lint image operator-image agent-image push
VERSION := $(shell ./build/git-version.sh)
RELEASE_VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse HEAD)

ifneq ($(VERSION), $(RELEASE_VERSION))
    VERSION := $(RELEASE_VERSION)-$(VERSION)
endif

REPO=github.com/pantheon-systems/container-linux-update-operator
# ko can't pass -ldflags any other way
GOFLAGS := -ldflags=-w
GOFLAGS := $(GOFLAGS) -ldflags=-X=$(REPO)/pkg/version.Version=$(RELEASE_VERSION)
GOFLAGS := "$(GOFLAGS) -ldflags=-X=$(REPO)/pkg/version.Commit=$(COMMIT)"

IMAGE_REPO ?= quay.io/getpantheon/container-linux-update-operator

KUBE_NAMESPACE ?= $(shell kubectl config get-contexts \
    | grep $(kubectl config current-context) | awk '{ print $NF}')

all: bin/* test lint coverage

test: deps
bin/*: deps

tools: export GO111MODULE=off
tools:
	go get -u "github.com/golangci/golangci-lint/cmd/golangci-lint" > /dev/null
	go get -u "github.com/ory/go-acc" > /dev/null

deps: tools
	go get ./...

export CGO_ENABLED := 0

# support locally executable operator binary for use with `kubectl proxy`
GOOS ?= $(shell go env GOOS)
bin/update-operator: test
	GOFLAGS=$(GOFLAGS) GOARCH=amd64 GOOS=$(GOOS) go build -o bin/update-operator \
        $(REPO)/cmd/update-operator

# default to linux because this binary is meant to only run on a linux host
GOOS ?= linux
bin/update-agent: test
	GOFLAGS=$(GOFLAGS) GOARCH=amd64 GOOS=$(GOOS) go build -o bin/update-agent \
        $(REPO)/cmd/update-agent

test:
	go test -v $(REPO)/pkg/...

lint:
	golangci-lint run

coverage:
	go-acc ./...

agent-image: bin/update-agent
	@docker build -t $(IMAGE_REPO):$(VERSION) --build-arg=cmd=update-agent .

operator-image: bin/update-operator
	@docker build -t $(IMAGE_REPO):$(VERSION) --build-arg=cmd=update-operator .

image: agent-image
image: operator-image

push: image
	@docker push $(IMAGE_REPO):$(VERSION)

ko:
	@ko apply -f k8s/daemonset.yaml

clean:
	rm -rf bin
