##############################################################
# This file is used by brazil-build through BMG.
##############################################################
include build/private/bgo_exports.makefile
include ${BGO_MAKEFILE}

pre-build:: imports-check-no-vendor fmt-check vet lint

.PHONY: lint
lint: lint-golangci

.PHONY: lint-golangci
lint-golangci:
	GO111MODULE=on golangci-lint run

.PHONY: imports-check-no-vendor
imports-check-no-vendor::
	$(eval DIFFS := $(shell goimports -l pkg cmd))
	$(at)if [ -n "$(DIFFS)" ]; then echo "Imports incorrectly formatted/ordered."; echo "Incorrectly formatted files: $(DIFFS)"; exit 1; fi
