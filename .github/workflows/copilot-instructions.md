# Authentication Service - Code Review Guidelines

This repository contains a Go authentication service for a Hotel Management System.

## General Review Rules

Review every Pull Request for:

- Correctness
- Security
- Maintainability
- Performance
- Error handling
- Test coverage
- Go best practices
- API design
- Database safety

Do not suggest changes merely for stylistic preference.

Only report issues that are meaningful and actionable.

## Go

Check for:

- Idiomatic Go
- Proper error handling
- Unnecessary complexity
- Incorrect pointer usage
- Goroutine leaks
- Data races
- Context propagation
- Resource leaks
- Improper use of interfaces
- Poor package boundaries
- Unnecessary allocations

Prefer simple and idiomatic Go solutions.

## Authentication Security

Pay special attention to:

- Password hashing
- Password verification
- JWT generation
- JWT validation
- Token expiration
- Token refresh
- Token revocation
- Authentication bypasses
- Authorization bypasses
- Role validation
- Privilege escalation
- Sensitive information leakage
- Secrets accidentally committed to source code
- Insecure logging of credentials or tokens

Passwords must never be stored or logged in plaintext.

JWT secrets and other credentials must never be hardcoded.

## API Security

Check for:

- Missing authentication middleware
- Missing authorization checks
- Improper input validation
- User ID manipulation
- IDOR vulnerabilities
- Excessive information disclosure
- Improper HTTP status codes
- Missing rate limiting where appropriate

## Database

Check for:

- SQL injection
- Unsafe queries
- Missing transactions
- Incorrect transaction handling
- Connection/resource leaks
- N+1 queries
- Missing indexes where clearly necessary

## Tests

Check whether new functionality has appropriate tests.

For authentication functionality, consider:

- Successful login
- Invalid credentials
- Invalid/expired tokens
- Unauthorized requests
- Forbidden requests
- Password change
- Profile updates
- Invalid input
- Security-sensitive edge cases

Do not require tests for trivial changes when tests would provide no meaningful value.

## Review Output

Prioritize findings by severity:

1. Critical security or correctness issues
2. High-impact bugs
3. Important maintainability issues
4. Minor improvements

Do not report formatting issues that are already handled by gofmt or golangci-lint unless there is a meaningful reason.

Do not approve code simply because the code compiles.

Focus on genuine problems rather than generating a large number of low-value comments.