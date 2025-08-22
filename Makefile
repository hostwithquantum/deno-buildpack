sample:=deno

buildpack?=.
bin_dir:=$(CURDIR)/bin
builder?=r.planetary-quantum.com/runway-public/runway-buildpack-stack:jammy-full

.PHONY: build
build:
	goreleaser build --single-target --clean --snapshot

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
	act -P ubuntu-latest=catthehacker/ubuntu:act-latest "pull_request"

.PHONY: smoke-%
smoke-%: test=$*
smoke-%:
	pack \
		build \
		test-$(test)-app \
		--builder $(builder) \
		--path ./samples/$(test) \
		--env "BP_LOG_LEVEL=DEBUG" \
		--pull-policy never \
		--buildpack $(buildpack)
