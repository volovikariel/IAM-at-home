# Gateway API

## POST /v1/users
### Description
Creates a new user

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

## GET /v1/users/{username}
### Description
Gets user info


### Responses
- 200 OK
```json
{ "username": "foo", "_rel": { "self": "/v1/users/foo", "session": "/v1/users/sessions/foo" } }
```
- 406 Not Acceptable

---

## PATCH /v1/users/{username}
### Description
Updates user

### Request
#### Content Type
application/json

#### Parameters
- password (`string`): ...
- session-token (`string`): ...

### Responses
- 200 OK
```json
{ "username": "foo", "_rel": { "self": "/v1/users/foo" } }
```
- 400 Bad Request
- 401 Unauthorized
- 403 Forbidden
- 415 Unsupported Media Type
- 422 Unprocessable Entity

---

## DELETE /v1/users/{username}
### Description
Deletes a user

### Request
#### Content Type
application/json

#### Parameters
- session-token (`string`): ...

### Responses
- 204 No Content
- 401 Unauthorized
- 403 Forbidden
- 400 Bad Request
- 415 Unsupported Media Type
- 422 Unprocessable Entity

---

## POST /v1/users/sessions
### Description
Starts a user session

### Request
#### Content Type
application/json

#### Parameters
- username (`string`): ...
- password (`string`): ...

### Responses
- 201 Created
```json
{ "session-token": "foo", "_rel": { "self": "/v1/users/sessions/foo" } }
```
- 400 Bad Request
- 401 Unauthorized
- 415 Unsupported Media Type
- 422 Unprocessable Entity

---

## DELETE /v1/users/sessions/{username}
### Description
Terminates session

### Request
#### Content Type
application/json

#### Parameters
- session-token (`string`): ...

### Responses
- 204 No Content
- 400 Bad Request
- 401 Unauthorized
- 415 Unsupported Media Type
- 422 Unprocessable Entity

---

##  other
### Description
Whenever the path does not have a resource associated with it, or when the method does not have a resource associated with it (e.g: user not found), or when the uri is just flatout wrong


### Responses
- 404 Not Found

---

