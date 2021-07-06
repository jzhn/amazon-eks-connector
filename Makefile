include build/private/bgo_exports.makefile
include ${BGO_MAKEFILE}

##############################################################
# Below targets are for dev/CI only.
# They are not executed by brazil-build
##############################################################
GORELEASER := $(shell command -v goreleaser 2> /dev/null)
MOCKERY := $(shell command -v mockery 2> /dev/null)

.PHONY: ci-release
# Turn on go modules - it's disabled in BMG.
ci-release:: GO111MODULE=on
ci-release::
ifndef GORELEASER
	$(error "goreleaser not found (`follow https://goreleaser.com/install/` to fix)")
endif
	$(GORELEASER) --skip-publish --rm-dist --snapshot

mocks-clean:
	rm -f ./pkg/**/mock_*.go

mocks-gen: mocks-clean
ifndef MOCKERY
	$(error "mockery not found (`follow https://github.com/vektra/mockery#installation` to fix)")
endif
	$(MOCKERY) \
      --case underscore \
      --all \
      --dir pkg \
      --inpackage
