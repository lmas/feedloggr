COVER := ".cover"

# Show available recipes by default
default:
    @just --list

# Run tests and log the test coverage
test:
    go test -v -coverprofile="{{COVER}}.out" ./...

# Run benchmarking
bench:
    go test -cover -bench=.

# Generate pretty coverage report
cover:
    go tool cover -html="{{COVER}}.out" -o="{{COVER}}.html"
    firefox "{{COVER}}.html"

# Runs source code linters (ignoring errors from gosec and govulncheck)
lint:
    go vet ./...
    - gosec -quiet -fmt=golint ./...
    - govulncheck ./...

# Updates 3rd party packages and tools
deps:
    go get -u ./...
    go mod tidy
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest

# Show documentation of public parts of package, in the current dir
[no-cd]
docs:
    go doc -all

# Builds the binary, with debug symbol table and DWARF gen disabled for smaller bin
build:
    go build -ldflags "-s -w"

# Clean up built binary and other temporary files (ignores errors from rm)
clean:
    go clean
    -rm {{COVER}}.*

