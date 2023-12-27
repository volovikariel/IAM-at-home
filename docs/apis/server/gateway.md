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
            <td>{username: foo, password: ...}</td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>201 Created</td>
                            <td>
                            {
                                "username": foo,
                                "_rel": {
                                    "self": "/v1/users/foo"
                                }
                            }</td>
                        </tr>
                        <tr><td>409 Conflict</td><td></td></tr>
                        <tr><td>400 Bad Request</td><td></td></tr>
                        <tr><td>415 Unsupported Media Type</td><td></td></tr>
                        <tr><td>422 Unprocessable Entity</td><td></td></tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>/v1/users/{username}</td>
            <td>GET</td>
            <td>Getting user info</td>
            <td></td>
            <td></td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>200 OK</td>
                            <td>
                            {
                                "username": foo,
                                "_rel": {
                                    "self": "/v1/users/foo",
                                    "session": "/v1/users/sessions/foo"
                                }
                            }</td>
                        </tr>
                        <tr><td>406 Not Acceptable</td><td></td></tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>/v1/users/{username}</td>
            <td>PATCH</td>
            <td>Update user</td>
            <td>application/json</td>
            <td>{password: ..., session-token: ...}</td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>200 OK</td>
                            <td>
                            {
                                "username": foo,
                                "_rel": {
                                    "self": "/v1/users/foo"
                                }
                            }</td>
                        </tr>
                        <tr><td>401 Unauthorized</td><td></td></tr>
                        <tr><td>403 Forbidden</td><td></td></tr>
                        <tr><td>400 Bad Request</td><td></td></tr>
                        <tr><td>415 Unsupported Media Type</td><td></td></tr>
                        <tr><td>422 Unprocessable Entity</td><td></td></tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>/v1/users/{username}</td>
            <td>DELETE</td>
            <td>Delete a user</td>
            <td>application/json</td>
            <td>{session-token: ...}</td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr><td>204 No Content</td><td></td></tr>
                        <tr><td>401 Unauthorized</td><td></td></tr>
                        <tr><td>403 Forbidden</td><td></td></tr>
                        <tr><td>400 Bad Request</td><td></td></tr>
                        <tr><td>415 Unsupported Media Type</td><td></td></tr>
                        <tr><td>422 Unprocessable Entity</td><td></td></tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>/v1/users/sessions</td>
            <td>POST</td>
            <td>Starting a session</td>
            <td>application/json</td>
            <td>{username: ..., password: ...}</td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>201 Created</td>
                            <td>
                            {
                                "session-token": foo,
                                "_rel": {
                                    "self": "/v1/users/sessions/foo"
                                }
                            }</td>
                        </tr>
                        <tr><td>401 Unauthorized</td><td></td></tr>
                        <tr><td>400 Bad Request</td><td></td></tr>
                        <tr><td>415 Unsupported Media Type</td><td></td></tr>
                        <tr><td>422 Unprocessable Entity</td><td></td></tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>/v1/users/sessions/{username}</td>
            <td>DELETE</td>
            <td>Terminate session</td>
            <td>application/json</td>
            <td>{session-token: ...}</td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr><td>204 No Content</td><td></td></tr>
                        <tr><td>401 Unauthorized</td><td></td></tr>
                        <tr><td>400 Bad Request</td><td></td></tr>
                        <tr><td>415 Unsupported Media Type</td><td></td></tr>
                        <tr><td>422 Unprocessable Entity</td><td></td></tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>other</td>
            <td></td>
            <td></td>
            <td></td>
            <td></td>
            <td>
                <table>
                    <thead>
                        <tr>
                            <th>Response Status</th>
                            <th>Response Object</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr><td>404 Not Found</td><td></td></tr>
                </table>
            </td>
        </tr>
    </tbody>
</table>