# Integration Tests

This directory contains integration tests for the Rojgarsetu backend API.

## Running Tests

Run tests that do not require PostgreSQL (the CI short-mode suite):
```bash
go test ./... -short -v -timeout=60s
```

Run the full suite, including database-backed tests, against a prepared PostgreSQL database:
```bash
go test ./tests/... -v -timeout=60s
```

Run specific test file:
```bash
go test ./tests/auth_test.go -v
```

Run with coverage:
```bash
go test ./tests/... -v -cover
```

## Test Files

- `auth_test.go`: Authentication and authorization tests
- `jobs_test.go`: Job-related endpoint tests
- `crawler_test.go`: Crawler functionality tests

## Test Coverage Target

Target: 80%+ coverage on handlers and services

## CI Integration

CI runs the short-mode suite without requiring a PostgreSQL service. Database-backed integration tests explicitly skip in short mode; run the full suite against a disposable, migrated test database before release. A CI pass is not evidence that the database integration tests ran.