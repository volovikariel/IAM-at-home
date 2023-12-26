### Table format
<table>
    <thead>
        <tr>
            <th rowspan="2">Endpoint</th>
            <th rowspan="2">Method</th>
            <th rowspan="2">Purpose</th>
            <th colspan="2">Accepts</th>
            <th rowspan="2">Responses</th>
        </tr>
        <tr>
            <th>Content-Type</th>
            <th>Example</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>/v1/users</td>
            <td>POST</td>
            <td>Creating a user</td>
            <td>application/json</td>
            <td>{username: ..., password: ...}</td>
            <td>
                <ul>
                    <li>201 Created</li>
                    <li>409 Conflict</li>
                    <li>400 Bad Request</li>
                    <li>415 Unsupported Media Type</li>
                    <li>422 Unprocessable Entity</li>
                </ul> 
            </td>
        </tr>
        <tr>
            <td>/v1/users/{username}</td>
            <td>GET</td>
            <td>Getting user info</td>
            <td></td>
            <td></td>
            <td>
                <ul>
                    <li>200 OK</li>
                    <li>406 Not Acceptable</li>
                </ul>
            </td>
        </tr>
        <tr>
            <td>/v1/users/{username}</td>
            <td>PATCH</td>
            <td>Update user</td>
            <td>application/json</td>
            <td>{password: ..., session-token: ...}</td>
            <td>
                <ul>
                    <li style="color: yellow;">204 No Content</li>
                    <li>401 Unauthorized</li>
                    <li>403 Forbidden</li>
                    <li>400 Bad Request</li>
                    <li>415 Unsupported Media Type</li>
                    <li>422 Unprocessable Entity</li>
                </ul>
            </td>
        </tr>
        <tr>
            <td>/v1/users/{username}</td>
            <td>DELETE</td>
            <td>Delete a user</td>
            <td>application/json</td>
            <td>{session-token: ...}</td>
            <td>
                <ul>
                    <li>204 No Content</li>
                    <li>401 Unauthorized</li>
                    <li>403 Forbidden</li>
                    <li>400 Bad Request</li>
                    <li>415 Unsupported Media Type</li>
                    <li>422 Unprocessable Entity</li>
                </ul>
            </td>
        </tr>
        <tr>
            <td>/v1/users/sessions</td>
            <td>POST</td>
            <td>Starting a session</td>
            <td>application/json</td>
            <td>{username: ..., password: ...}</td>
            <td>
                <ul>
                    <li>201 Created</li>
                    <li>401 Unauthorized</li>
                    <li>400 Bad Request</li>
                    <li>415 Unsupported Media Type</li>
                    <li>422 Unprocessable Entity</li>
                </ul> 
            </td>
        </tr>
        <tr>
            <td>/v1/users/sessions/{username}</td>
            <td>DELETE</td>
            <td>Terminate session</td>
            <td>application/json</td>
            <td>{session-token: ...}</td>
            <td>
                <ul>
                    <li>204 No Content</li>
                    <li>401 Unauthorized</li>
                    <li>400 Bad Request</li>
                    <li>415 Unsupported Media Type</li>
                    <li>422 Unprocessable Entity</li>
                </ul>
            </td>
        </tr>
        <tr>
            <td>other</td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td>
                <ul>
                    <li>404 Not found</li>
                </ul>
            </td>
        </tr>
    </tbody>
</table>

#### Text version

### Registering
#### Endpoint: /v1/users

##### Method: POST

Purpose: Creating a user

Content-Type: application/json

Example: {username: ..., password: ...}

Responses:

201 Created

409 Conflict

400 Bad Request

415 Unsupported Media Type

422 Unprocessable Entity

### Endpoint: /v1/users/{username}
#### Method: GET

Purpose: get user

Content-Type: application/json

Responses:

200 OK

406 Not Acceptable

#### Method: PATCH 

Purpose: update user

Content-Type: application/json

Example: {password: ..., session-token: ...}

Responses:

204 No Content

401 Unauthorized

403 Forbidden

400 Bad Request

415 Unsupported Media Type

422 Unprocessable Entity


##### Method: DELETE

Purpose: delete user

Content-Type: application/json

Example: {session-token: ...}

Responses:

204 No Content

401 Unauthorized

403 Forbidden

400 Bad Request

415 Unsupported Media Type

422 Unprocessable Entity

### Logging in/out & authenticating
#### Endpoint: /v1/users/sessions

##### Method: POST

Purpose: create session

Content-Type: application/json

Example: {username: ..., password: ...}

Responses:
201 Created

401 Unauthorized

400 Bad Request

415 Unsupported Media Type

422 Unprocessable Entity


#### Endpoint: /v1/users/sessions/{username}

##### Method: DELETE

Purpose: terminate session

Content-Type: application/json

Example: {session-token: ...}

Responses:

204 No Content

401 Unauthorized

400 Bad Request

415 Unsupported Media Type

422 Unprocessable Entity

All non-specified URIs return a 404 Not Found.