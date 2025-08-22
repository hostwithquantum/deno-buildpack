sample:=deno

buildpack?=.
builder?=r.planetary-quantum.com/runway-public/runway-buildpack-stack:jammy-full

BUILD_DIR:=./build
VERSION?=dev
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: setup
setup:
	pack config default-builder $(builder)
	pack config trusted-builders add $(builder)

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
		--buildpack $(buildpack)

.PHONY: prep
prep:
	mkdir -p $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/linux/amd64/bin
	mkdir -p $(BUILD_DIR)/linux/arm64/bin

	cp dist/build_linux_amd64*/build $(BUILD_DIR)/linux/amd64/bin/
	cp dist/detect_linux_amd64*/detect $(BUILD_DIR)/linux/amd64/bin/
	cp dist/build_linux_arm64*/build $(BUILD_DIR)/linux/arm64/bin/
	cp dist/detect_linux_arm64*/detect $(BUILD_DIR)/linux/arm64/bin/

	cp LICENSE README.md buildpack.toml package.toml $(BUILD_DIR)/
	sed -i.bak -E "s/__replace__/$(VERSION)/" $(BUILD_DIR)/buildpack.toml
	rm -f $(BUILD_DIR)/buildpack.toml.bak