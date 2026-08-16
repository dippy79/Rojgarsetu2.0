# Integration Tests

This directory contains integration tests for the Rojgarsetu backend API.

## Running Tests

Run all tests:
```bash
cd backend_go
go test ./tests/... -v
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

These tests are automatically run in the GitHub Actions CI/CD pipeline on every push to master.