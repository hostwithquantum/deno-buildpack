sample:=deno

buildpack?=.
bin_dir:=$(CURDIR)/bin
builder?=r.planetary-quantum.com/runway-public/paketo-builder:full

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

.PHONY: smoke
smoke:
	pack \
		build \
		test-$(sample)-app \
		--path ./samples/$(sample) \
		--buildpack $(buildpack)
