![Workflow Badge](https://github.com/volovikariel/IdentityManager/actions/workflows/go.yml/badge.svg)

# Table of contents
1. [Running individual components](#running-individual-components)
   1. [Prerequisites](#running-individual-components-prerequisites)
   1. [Gateway API Server](#running-individual-components-gateway-api-server)
1. [Running a cluster of components](#running-component-clusters) (Allows for manual/automatic scaling)
   1. [Prerequisites](#running-component-clusters-prerequisites)
   1. [Gateway API Server](#running-component-clusters-gateway-api-server)
1. [Load Testing](#load-tests)
   1. [Gateway API Server](#load-testing-gateway-api-server)
1. [Documentation](#documentation)
   1. [API](#documentation-api)
      1. [Gateway API Docs](#documentation-api-gateway-api-server)
   1. [Architecture diagrams](#documentation-architecture-diagrams)
   
<a name="running-individual-components"></a>
# Running individual components

<a name="running-individual-components-prerequisites"/></a>
## Prerequisites: Docker

<a name="running-individual-components-gateway-api-server"/></a>
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

<a name="running-component-clusters"/></a>
# Running component clusters
<a name="running-component-clusters-prerequisites"/></a>
## Prerequisites
1. Docker installed (to build the image).
1. Minikube installed (to run a Kubernetes cluster on your machine).
1. Kubectl installed (to manage your cluster's deployments, services, etc.); alternatively - you can run the following and use Minikube's built in Kubectl: `alias kubectl="minikube kubectl --"`

As we're using local images in our clusters, do the following before building the images:
```bash
eval $(minikube docker-env)
```
Alternatively, if the image is built, you can run
```bash
minikube image load gateway-api
```

<a name="running-component-clusters-gateway-api-server"/></a>
## Gateway API Server
Build its Docker image (see steps above).

Create a Kubernetes cluster:
```bash
minikube start
```

Create the Gateway API Server Deployment:
```bash
kubectl apply -f ./build/server/gateway/deployment.yaml
```

<a id="intro" name="intro"></a>
Create the Gateway API Server Service:
```bash
kubectl apply -f ./build/server/gateway/service.yaml
```

**Note**: You should now be able to access the Gateway API Server, as a test you can run:
```bash
PORT=$(kubectl get service gateway-api-service -o=jsonpath='{.spec.ports[0].nodePort}'); \
curl http://$(minikube ip):$PORT
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

<a name="load-testing"/></a>
# Load tests
Instructions on how to run load tests [here](/internal/docs/tests/load/README.md).

<a name="load-testing-gateway-api-server"/></a>
## Gateway API Server
[load test scenario](https://github.com/volovikariel/IdentityManager/blob/d87ba775da37ad427be70f47c55d64df7268eaaf/internal/docs/tests/load/gateway_api.js)):

[📈1 replica📈](https://volovikariel.github.io/IdentityManager/tests/load/gateway/1_replicas_report.html)

[📈2 replicas📈](https://volovikariel.github.io/IdentityManager/tests/load/gateway/2_replicas_report.html)

[📈6 replicas📈](https://volovikariel.github.io/IdentityManager/tests/load/gateway/6_replicas_report.html)

**Note**: Better tests are planned once the application is completed:
- Breakpoint testing: Very slowly scaling up #VUs (to ensure that your SLOs are met)
- Stress testing: Higher than expected average load for a medium length of time (see whether the system scales to adjust to it properly)
- Spike testing: Insane load for a very short amount of time (see whether the system recovers gracefully from any failures that may occur)

<a name="documentation"/></a>
# Documentation
<a name="documentation-api"/></a>
## APIs
<a name="documentation-api-gateway-api-server"/></a>
### Gateway API Server
[Docs](https://volovikariel.github.io/IdentityManager/apis/server/gateway_api.html)

<a name="documentation-architecture-diagrams"/></a>
## Diagrams
### Interactions
![User interactions](diagrams/user_interactions.svg)
![Admin interactions](diagrams/admin_interactions.svg)
### Services
![Registration service](diagrams/registration_service.svg)
