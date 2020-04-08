.PHONY:	all release-bin image clean test vendor
export CGO_ENABLED:=0

VERSION=$(shell ./build/git-version.sh | head -c 7)
RELEASE_VERSION=$(shell cat VERSION)
COMMIT=$(shell git rev-parse HEAD)

REPO=github.com/pantheon-systems/container-linux-update-operator
LD_FLAGS="-w -X $(REPO)/pkg/version.Version=$(RELEASE_VERSION) -X $(REPO)/pkg/version.Commit=$(COMMIT)"

IMAGE_REPO?=quay.io/getpantheon/container-linux-update-operator

all: bin/update-agent bin/update-operator

GOOS ?= $(shell go env GOOS)
bin/update-operator:
	GOARCH=amd64 GOOS=$(GOOS) go build -o bin/update-operator -ldflags $(LD_FLAGS) \
      $(REPO)/cmd/update-operator

bin/update-agent:export GOOS=linux
bin/update-agent:
	GOARCH=amd64 GOOS=$(GOOS) go build -o bin/update-agent -ldflags $(LD_FLAGS) \
      $(REPO)/cmd/update-agent

release-bin:
	./build/build-release.sh

test:
	go test -v $(REPO)/pkg/...

image: release-bin
	@sudo docker build --rm=true -t $(IMAGE_REPO):$(VERSION) .

docker-push: image
	@sudo docker push $(IMAGE_REPO):$(VERSION)

vendor:
	glide update --strip-vendor
	glide-vc --use-lock-file --no-tests --only-code

clean:
	rm -rf bin
