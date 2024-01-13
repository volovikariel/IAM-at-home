![Workflow Badge](https://github.com/volovikariel/IdentityManager/actions/workflows/go.yml/badge.svg)

# Running the standalone parts
## Gateway API Server
Build the image: `docker build -f build/server/gateway/Dockerfile -t gateway-api .`

Run the image: `HOST_NAME='' HOST_PORT=10000 CONTAINER_PORT=8080; docker run -d --expose "$CONTAINER_PORT" -p "$HOST_PORT:$CONTAINER_PORT" --name gateway-api gateway-api:latest -h "$HOST_NAME" -p "$CONTAINER_PORT"`

**Note**: You can specify the host name, host port, and container ports in the `docker run` command.

You should now be able to access the Gateway API Server at: `http://localhost:HOST_PORT` (e.g: `curl localhost:10000`)

Stop the container: `docker stop gateway-api`
Remove the container: `docker rm gateway-api`

# APIs
## Gateway API Server
[Docs](https://volovikariel.github.io/IdentityManager/apis/server/gateway_api.html)

# Diagrams
## Interactions
![User interactions](diagrams/user_interactions.svg)
![Admin interactions](diagrams/admin_interactions.svg)
## Services
![Registration service](diagrams/registration_service.svg)
