sample:=deno

buildpack?=.
bin_dir:=$(CURDIR)/bin
builder:=r.planetary-quantum.com/runway-public/paketo-builder:full

.PHONY: build
build: clean
	GOOS=linux go build \
		-ldflags="-s -w" \
		-o "$(bin_dir)/detect" \
		"$(CURDIR)/cmd/detect/main.go"

	GOOS=linux go build \
		-ldflags="-s -w" \
		-o "$(bin_dir)/build" \
		"$(CURDIR)/cmd/build/main.go"

.PHONY: clean
clean:
	rm -f $(bin_dir)/detect
	rm -f $(bin_dir)/run


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
		--buildpack $(buildpack)
