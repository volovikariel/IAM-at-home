We use [Grafana K6](https://github.com/grafana/k6) to run our load tests.

To visualize the load test results, we use [xk6 dashboard](https://github.com/grafana/xk6-dashboard).

Here's an example of performing a load test and running it on our Kubernetes cluster:

```js
// tests/load/gateway_api.js
const HOST = __ENV.HOST;
const PORT = __ENV.PORT;

export const options = {
  scenarios: {
    ramping_unsupported_path: {
      executor: "ramping-vus",
      startVus: 1,
      exec: "getPath",
      env: {
        path: "unsupported-path",
      },
      stages: [
        { duration: "5s", target: 10 },
        { duration: "5s", target: 100 },
        { duration: "5s", target: 1000 },
        { duration: "5s", target: 2500 },
        { duration: "5s", target: 5000 },
        { duration: "5s", target: 7500 },
        { duration: "5s", target: 10000 },
      ],
    },
    // ...
  }
};

export function getPath() {
  const url = `http://${HOST}:${PORT}/${__ENV.path}`;
  http.get(url);
  sleep(1);
}
```

Which can then be run with the following
```bash
NUM_REPLICAS=$(kubectl get deployments gateway-api-deployment -o jsonpath='{.spec.replicas}'); \
HOST=$(minikube ip) \
PORT=30000 \
k6 run \
--out web-dashboard=export="docs/tests/load/gateway/${NUM_REPLICAS}_replicas_report.html" \
internal/docs/tests/load/gateway_api.js
```

**NOTE**: The `HOST` and `PORT` can be edited to match your environment.

**NOTE**: To use `--out web-dashboard` you need to have Grafana Xk6 Dashboard installed, see [here](https://github.com/grafana/xk6-dashboard).