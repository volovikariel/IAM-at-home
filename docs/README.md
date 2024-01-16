![Workflow Badge](https://github.com/volovikariel/IdentityManager/actions/workflows/go.yml/badge.svg)

# Running the standalone parts (Docker)
## Gateway API Server
Build the image: 
```bash
docker build -f build/server/gateway/Dockerfile -t gateway-api .
```

Run the image:
```bash
HOST_NAME='' HOST_PORT=10000 CONTAINER_PORT=8080; \
docker run -d \
--expose "$CONTAINER_PORT" \
-p "$HOST_PORT:$CONTAINER_PORT" \
--name gateway-api \
gateway-api:latest \
-h "$HOST_NAME" \
-p "$CONTAINER_PORT"
```

**Note**: You can specify the host name, host port, and container ports in the `docker run` command.

You should now be able to access the Gateway API Server at: `http://localhost:$HOST_PORT` (e.g: `curl localhost:$HOST_PORT`, or `curl localhost:10000`)

Stop the container:
```bash
docker stop gateway-api
```

Remove the container:
```bash
docker rm gateway-api
```

# Running the individual parts (Kubernetes)
As we'll be using a local image and deploying it to a Kubernetes cluster, do the following before building the docker image:
```bash
eval $(minikube docker-env)
```

## Gateway API Server
Build its Docker image.

Create a Kubernetes cluster:
```bash
minikube start
```

Create the Gateway API Server Deployment:
```bash
kubectl apply -f ./build/server/gateway/deployment.yaml
```

Create the Gateway API Server Service:
```bash
kubectl apply -f ./build/server/gateway/service.yaml
```

**Note**: You should now be able to access the Gateway API Server at `http://$(minikube ip):30000` (e.g: `curl http://$(minikube ip):30000`)

Delete the service:
```bash
kubectl delete service gateway-api-service
```

Delete the deployment:
```bash
kubectl delete deployment gateway-api-deployment
```

Delete the cluster:
```bash
minikube delete
```

To scale the number of replicas after having deployed the Gateway API Server:
```bash
REPLICAS=2; \
kubectl scale -f ./build/server/gateway/deployment.yaml --replicas=$REPLICAS
```

Alternatively, you can modify the `replicas` field in the `deployment.yaml` file, then run:
```bash 
kubectl apply -f ./build/server/gateway/deployment.yaml
```

# APIs
## Gateway API Server
[Docs](https://volovikariel.github.io/IdentityManager/apis/server/gateway_api.html)

# Diagrams
## Interactions
![User interactions](diagrams/user_interactions.svg)
![Admin interactions](diagrams/admin_interactions.svg)
## Services
![Registration service](diagrams/registration_service.svg)
