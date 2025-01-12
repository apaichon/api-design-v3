# User Management API Documentation

## Overview

The User Management API allows you to create, read, update, and delete user accounts in your application. This RESTful API uses JSON for request and response payloads and follows standard HTTP methods and status codes.

## Base URL

```
https://api.example.com/v1
```

## Authentication

All API requests require authentication using an API key that should be included in the HTTP header:

```
Authorization: Bearer YOUR_API_KEY
```

### API Key Format and Validation

API keys follow this format: `exkey_live_xxxxxxxxxxxxxxxxxxxxxxxxxxxx`

- Prefix: `exkey_` indicates an API key
- Environment: `live_` for production or `test_` for testing
- Key: 32 character string containing lowercase letters and numbers
- Total length: 44 characters

#### Validation Rules

1. Format validation (regex):
```
^exkey_(live|test)_[a-z0-9]{32}$
```

2. Security requirements:
   - Keys must be transmitted only over HTTPS
   - Keys must be stored securely using encryption at rest
   - Keys must never be logged or exposed in client-side code
   - Each service/application should use a unique API key

#### API Key Lifecycle

1. Generation
   - Keys are generated through the developer dashboard
   - Each key has associated metadata (creation date, last used, permissions)
   - Maximum of 10 active keys per account

2. Validation Process
   ```
   - Check key format
   - Verify key exists in database
   - Check key is not expired or revoked
   - Validate permissions for requested endpoint
   - Update last_used timestamp
   ```

3. Error Handling
   - Invalid format: HTTP 401 with `INVALID_KEY_FORMAT`
   - Invalid key: HTTP 401 with `INVALID_API_KEY`
   - Expired key: HTTP 401 with `EXPIRED_API_KEY`
   - Insufficient permissions: HTTP 403 with `INSUFFICIENT_PERMISSIONS`

4. Security Measures
   - Rate limiting per API key
   - Automatic key rotation every 90 days
   - Instant key revocation capability
   - Audit logging of key usage

