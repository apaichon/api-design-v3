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

## Rate Limiting

- Rate limit: 1000 requests per hour per API key
- Exceeded requests receive HTTP 429 (Too Many Requests)
- Rate limit headers included in all responses:
  - X-RateLimit-Limit: Total requests allowed per hour
  - X-RateLimit-Remaining: Remaining requests for the current hour
  - X-RateLimit-Reset: Unix timestamp when the limit resets

## Endpoints

### Create User

Creates a new user account.

**POST** `/users`

#### Request Body

```json
{
  "email": "string",
  "username": "string",
  "firstName": "string",
  "lastName": "string",
  "password": "string"
}
```

#### Parameters

| Field      | Type   | Required | Description                              |
|------------|--------|----------|------------------------------------------|
| email      | string | Yes      | User's email address                     |
| username   | string | Yes      | Unique username (3-30 characters)        |
| firstName  | string | Yes      | User's first name                        |
| lastName   | string | Yes      | User's last name                         |
| password   | string | Yes      | Password (min 8 characters)              |

#### Response

**201 Created**

```json
{
  "id": "string",
  "email": "string",
  "username": "string",
  "firstName": "string",
  "lastName": "string",
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)"
}
```

#### Error Responses

| Status Code | Description                                          |
|------------|------------------------------------------------------|
| 400        | Bad Request - Invalid parameters                      |
| 409        | Conflict - Email or username already exists           |
| 429        | Too Many Requests - Rate limit exceeded              |

### Get User

Retrieves user details by ID.

**GET** `/users/{userId}`

#### Parameters

| Parameter | Type   | Location | Required | Description     |
|-----------|--------|----------|----------|-----------------|
| userId    | string | path     | Yes      | Unique user ID |

#### Response

**200 OK**

```json
{
  "id": "string",
  "email": "string",
  "username": "string",
  "firstName": "string",
  "lastName": "string",
  "createdAt": "string (ISO 8601)",
  "updatedAt": "string (ISO 8601)"
}
```

#### Error Responses

| Status Code | Description                             |
|-------------|-----------------------------------------|
| 404         | Not Found - User ID doesn't exist       |
| 429         | Too Many Requests - Rate limit exceeded |

### Update User

Updates an existing user's information.

**PUT** `/users/{userId}`

#### Parameters

| Parameter | Type   | Location | Required | Description     |
|-----------|--------|----------|----------|-----------------|
| userId    | string | path     | Yes      | Unique user ID |

#### Request Body

```json
{
  "email": "string",
  "firstName": "string",
  "lastName": "string"
}
```

#### Response

**200 OK**

```json
{
  "id": "string",
  "email": "string",
  "username": "string",
  "firstName": "string",
  "lastName": "string",
  "updatedAt": "string (ISO 8601)"
}
```

#### Error Responses

| Status Code | Description                             |
|-------------|-----------------------------------------|
| 400         | Bad Request - Invalid parameters        |
| 404         | Not Found - User ID doesn't exist       |
| 409         | Conflict - Email already exists         |
| 429         | Too Many Requests - Rate limit exceeded |

### Delete User

Deletes a user account.

**DELETE** `/users/{userId}`

#### Parameters

| Parameter | Type   | Location | Required | Description     |
|-----------|--------|----------|----------|-----------------|
| userId    | string | path     | Yes      | Unique user ID |

#### Response

**204 No Content**

#### Error Responses

| Status Code | Description                             |
|-------------|-----------------------------------------|
| 404         | Not Found - User ID doesn't exist       |
| 429         | Too Many Requests - Rate limit exceeded |

## Error Handling

All error responses follow this format:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
```

### Common Error Codes

| Code              | Description                                          |
|-------------------|------------------------------------------------------|
| INVALID_REQUEST   | Request validation failed                            |
| RESOURCE_EXISTS   | Resource already exists                              |
| RESOURCE_MISSING  | Resource not found                                   |
| RATE_LIMITED     | Rate limit exceeded                                  |

## Data Types

### User Object

| Field      | Type   | Description                              |
|------------|--------|------------------------------------------|
| id         | string | Unique identifier                        |
| email      | string | User's email address                     |
| username   | string | Unique username                          |
| firstName  | string | User's first name                        |
| lastName   | string | User's last name                         |
| createdAt  | string | ISO 8601 timestamp of creation          |
| updatedAt  | string | ISO 8601 timestamp of last update       |

## Best Practices

1. Always use HTTPS for API requests
2. Store API keys securely and never expose them in client-side code
3. Implement proper error handling for all possible response codes
4. Use appropriate HTTP methods for different operations
5. Follow rate limiting guidelines to avoid service disruption

## SDK Support

Official SDK libraries are available for:
- Python
- JavaScript
- Java
- Go
- Ruby

## Changelog

### v1.1.0 (2024-01-15)
- Added support for bulk user operations
- Improved rate limiting headers
- Added new error codes for validation

### v1.0.0 (2023-12-01)
- Initial API release
- Basic CRUD operations for users
- Authentication and rate limiting implementation

## Support

For technical support or questions about the API:
- Email: api-support@example.com
- Documentation: https://docs.example.com
- Status Page: https://status.example.com