sample:=deno

buildpack?=.
bin_dir:=$(CURDIR)/bin
builder?=r.planetary-quantum.com/runway-public/paketo-builder:full

.PHONY: build
build: clean bin/detect bin/build

bin/build:
	GOOS=linux go build \
		-ldflags="-s -w" \
		-o "$(bin_dir)/build" \
		"$(CURDIR)/cmd/build/main.go"

bin/detect:
	GOOS=linux go build \
		-ldflags="-s -w" \
		-o "$(bin_dir)/detect" \
		"$(CURDIR)/cmd/detect/main.go"

.PHONY: clean
clean:
	rm -f $(bin_dir)/detect
	rm -f $(bin_dir)/build


.PHONY: setup
setup:
	pack config default-builder $(builder)
	pack config trusted-builders add $(builder)


.PHONY: test
test: build
	pack \
		build \
		test-$(sample)-app \
		--path ./samples/$(sample) \
		--buildpack .


.PHONY: act-pr
act-pr:
	act "pull_request" \
		-s BP_QUANTUM_DOCKER_USERNAME \
		-s BP_QUANTUM_DOCKER_PASSWORD

.PHONY: smoke
smoke:
	pack \
		build \
		test-$(sample)-app \
		--path ./samples/$(sample) \
		--buildpack $(buildpack)
