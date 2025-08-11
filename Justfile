# -*- Justfile -*-

coverage_file := "coverage.out"

# List the available justfile recipes.
[group('general')]
@default:
  just --list --unsorted

# List the lines of code in the project.
[group('general')]
loc:
  scc --remap-unknown "-*- Justfile -*-":"justfile"

# Format and vet Go code. Runs before tests.
[group('test')]
check:
	go fmt ./...
	go vet ./...

# Lint using golangci-lint
[group('test')]
lint:
  golangci-lint run --config .golangci.yaml

# Run the unit tests.
[group('test')]
unit *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -short

# Run the integration tests.
[group('test')]
int *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -run Integration

# Run the end-to-end tests.
[group('test')]
e2e *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -run E2E

# HTML report for unit (default), int, e2e, or all tests.
[group('test')]
cover test='unit': check
  go test ./... -vet=off -coverprofile={{coverage_file}} \
  {{ if test == 'all' { '' } \
    else if test == 'int' { '-run Integration' } \
    else if test == 'e2e' { '-run E2E' } \
    else { '-short' } }}
  go tool cover -html={{coverage_file}}

# List the outdated direct dependencies (slow to run).
[group('dependencies')]
outdated:
  # (requires https://github.com/psampaz/go-mod-outdated).
  go list -u -m -json all | go-mod-outdated -update -direct

# Update the given module to the latest version.
[group('dependencies')]
update mod:
  go get -u {{mod}}
  go mod tidy

# Update all modules.
[group('dependencies')]
updateall:
  go get -u ./...
  go mod tidy

# Run go mod tidy and verify.
[group('dependencies')]
tidy:
  go mod tidy
  go mod verify
