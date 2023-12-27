# API Documentation

## POST /v1/users
### Description
Creating a new user

### Request
#### Content Type
application/json

#### Parameters
- username (`string`): foo
- password (`string`): ...

### Responses
- 201 Created
```json
{ "username": "foo", "_rel": { "self": "/v1/users/foo" } }
```
- 400 Bad Request
- 409 Conflict
- 415 Unsupported Media Type
- 422 Unprocessable Entity

---

