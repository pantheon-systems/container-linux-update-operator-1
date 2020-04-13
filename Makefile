.PHONY:	all image clean test coverage lint image operator-image agent-image push
VERSION := $(shell ./build/git-version.sh)
RELEASE_VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse HEAD)

ifneq ($(VERSION), $(RELEASE_VERSION))
    VERSION := $(RELEASE_VERSION)-$(VERSION)
endif

REPO=github.com/pantheon-systems/cos-update-operator
# ko can't pass -ldflags any other way
GOFLAGS := -ldflags=-w
GOFLAGS := $(GOFLAGS) -ldflags=-X=$(REPO)/pkg/version.Version=$(RELEASE_VERSION)
GOFLAGS := "$(GOFLAGS) -ldflags=-X=$(REPO)/pkg/version.Commit=$(COMMIT)"

OPERATOR_IMAGE_REPO ?= quay.io/getpantheon/cos-update-operator-operator
AGENT_IMAGE_REPO ?= quay.io/getpantheon/cos-update-operator-agent

KUBE_NAMESPACE ?= $(shell kubectl config get-contexts \
    | grep $(kubectl config current-context) | awk '{ print $NF}')

# set GOGC to mitigate OOMing and set lint cache location for use with circleci cache
ifdef CIRCLECI
    export GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache
    export GOGC=50
endif

all: bin/* test lint coverage

# used for caching golangci-lint data in circleci
master-sha:
	@git fetch origin && git rev-parse origin/master > master_sha

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
	golangci-lint run --verbose

coverage:
	go-acc ./...

agent-image: bin/update-agent
	@docker build -t $(AGENT_IMAGE_REPO):$(VERSION) --build-arg=cmd=update-agent .

operator-image: bin/update-operator
	@docker build -t $(OPERATOR_IMAGE_REPO):$(VERSION) --build-arg=cmd=update-operator .

image: agent-image
image: operator-image

push-agent: agent-image
	@docker push $(AGENT_IMAGE_REPO):$(VERSION)

push-operator: operator-image
	@docker push $(OPERATOR_IMAGE_REPO):$(VERSION)

push: push-agent push-operator

ko:
	@ko apply -f k8s/daemonset.yaml

clean:
	rm -rf bin
