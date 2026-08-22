---
name: authentication-security-review
description: Review authentication and authorization code for security vulnerabilities.
---

# Authentication Security Review

Review changes involving authentication, authorization, users, sessions,
passwords, JWTs, tokens, middleware, and permissions.

## Check

### Passwords

- Passwords must never be stored in plaintext.
- Password hashing must use an appropriate password hashing algorithm.
- Passwords must never appear in logs.
- Passwords must never appear in API responses.

### JWT

Check:

- Signature validation
- Algorithm validation
- Expiration validation
- Issuer validation where applicable
- Audience validation where applicable
- Token type validation
- Proper secret/key handling
- Token replay risks

Do not assume that decoding a JWT means it is valid.

### Authorization

Check every protected endpoint for:

- Authentication enforcement
- Role/permission validation
- Resource ownership validation
- Privilege escalation

Look specifically for IDOR vulnerabilities.

For example, an endpoint such as:

GET /users/{id}

must not automatically allow an authenticated user to retrieve another
user's private information.

### Input Validation

Look for:

- Missing validation
- Injection vulnerabilities
- Unexpected input types
- Excessively large input
- Malformed tokens
- Invalid user IDs

### Secrets

Flag:

- Hardcoded JWT secrets
- Database passwords
- API keys
- Private keys
- Credentials committed to source

Environment variables or an appropriate secret-management system should be used.

## Review Philosophy

Only report actual or strongly plausible security issues.

Do not report speculative vulnerabilities without explaining the attack path.