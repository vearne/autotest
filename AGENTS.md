# AGENTS.md

This file provides coding guidelines and project conventions for agentic developers working on this repository.

## Build & Test Commands

### Build
```bash
# Build the main binary
make build
# Or directly
go build -o autotest ./

# Cross-platform builds
make xgo  # Builds for linux/amd64, darwin/amd64, darwin/arm64
```

### Lint
```bash
# Run golangci-lint with project configuration
golangci-lint run
```

### Test
```bash
# Run all tests
go test ./...

# Run single test file
go test ./internal/rule -v

# Run specific test function
go test ./internal/rule -run TestXpathFind -v

# Run with coverage and race detection
go test -v -race -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out
```

### Run the Tool
```bash
# Validate config
autotest test --config-file=CONFIG_FILE

# Run tests
autotest run --config-file=CONFIG_FILE
autotest run --config-file=CONFIG_FILE --environment=dev

# Extract XPath
autotest extract --xpath=XPATH --json=JSON
```

## Code Style Guidelines

### Import Ordering
Imports MUST follow this 3-group structure with blank lines between groups:
1. Standard library imports (alphabetically sorted)
2. Third-party imports (alphabetically sorted)
3. Local/project imports (alphabetically sorted)

```go
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/antchfx/jsonquery"
	"github.com/urfave/cli/v3"

	"github.com/vearne/autotest/internal/resource"
	"github.com/vearne/autotest/internal/rule"
)
```

### Naming Conventions
- **Exported functions/variables**: PascalCase (`RunTestCases`, `ErrorIDduplicate`)
- **Unexported functions/variables**: camelCase (`convStr`, `confFilePath`)
- **Constants**: PascalCase or ALL_CAPS (`StateNotExecuted`, `ErrorIDduplicate`)
- **Interfaces**: PascalCase, descriptive names (`VerifyRule`, `IdItem`)
- **Method receivers**: Single lowercase letter (`func (r *Rule) Name()`)
- **Structs**: PascalCase (`UnifiedTestResults`, `HttpTestCaseResult`)

### Error Handling
- Check errors immediately after operations
- Log errors with context before returning
- Use `fmt.Errorf("...: %w", err)` to wrap errors when helpful
- Return errors, don't panic in application code

```go
err := resource.ParseConfigFile(confFilePath)
if err != nil {
	slog.Error("config file parse error, %v", err)
	return err
}
```

### Testing
- Place tests in `*_test.go` files alongside source code
- Use `github.com/stretchr/testify/assert` for assertions
- Use table-driven tests with struct slices for multiple test cases
- Test function naming: `func TestFunctionName(t *testing.T)`
- Subtests: `t.Run(tt.name, func(t *testing.T) { ... })`

### Package Organization
- `internal/` - Private packages (command, rule, util, config, model, resource, luavm)
- `example/` - Example applications
- `test/` - Test utilities
- `consts/` - Constants
- Root level - Main application entry point

### Comments & Documentation
- Mixed Chinese and English comments (Chinese for internal logic, English for exported APIs)
- Comments precede the item they document
- Simple style, no formal godoc format consistently required

```go
// UnifiedTestResults 统一的测试结果
type UnifiedTestResults struct {}

// CombineResults 合并HTTP和gRPC测试结果
func CombineResults(httpResults, grpcResults *UnifiedTestResults) *UnifiedTestResults
```

## Linter Configuration

The project uses golangci-lint with the following enabled linters:
- `copyloopvar` - Detects loop var copies
- `errcheck` - Checks for unchecked errors
- `govet` - Go vet static analysis
- `ineffassign` - Detects ineffectual assignments
- `staticcheck` - Go static analysis
- `unused` - Detects unused code

Disabled/Excluded:
- `fieldalignment` - Struct field alignment optimization
- Paths: `scripts`, `example`, `test`, third_party code
- Generated code: Lax handling

## Configuration Files

- `.golangci.yml` - Linting configuration
- `go.mod` - Go module definition (Go 1.24.0+)
- `Makefile` - Build and cross-platform compilation

## CI/CD

- Lint and build run on every push/PR to main branch
- Tests run only on releases via `.github/workflows/release.yml`
- No automated tests on main branch pushes/PRs

## Important Notes

- No Cursor or Copilot rules exist for this project
- The project prioritizes integration testing via the autotest tool over extensive unit tests
- Use `go test ./...` for unit tests, `autotest run` for integration tests
