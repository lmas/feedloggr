
COVER=.cover.out
COVER_HTML=.cover.html
# FUZZ_DIR=./fuzz
# FUZZ_BUILD=$(FUZZ_DIR)/fuzz.zip

.PHONY: test
test:
	go test -v -race -coverprofile="${COVER}" ./...

.PHONY: cover
cover:
	go tool cover -html="${COVER}" -o="${COVER_HTML}"

.PHONY: bench
bench:
	go test -race -cover -bench=.

.PHONY: lint
lint: govet golint gosec

.PHONY: govet
govet:
	go vet ./...

.PHONY: golint
golint:
	golint ./...

.PHONY: gosec
gosec:
	gosec -quiet -fmt=golint ./...

# .PHONY: fuzz
# fuzz:
	# go-fuzz-build -o "$(FUZZ_BUILD)" ./internal/
	# go-fuzz -bin "$(FUZZ_BUILD)" -workdir "$(FUZZ_DIR)"

.PHONY: build
build:
	go build -ldflags "-s -w" feedloggr.go

.PHONY: clean
clean:
	go clean
	rm "${COVER}" "${COVER_HTML}"
