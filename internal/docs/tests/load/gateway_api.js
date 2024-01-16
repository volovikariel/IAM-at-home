import http from "k6/http";
import { sleep, check } from "k6";

// If you're running in a Minikube cluster: result of `$ minikube ip`
// If you're running in a Docker container: localhost
const HOST = __ENV.HOST || "localhost";
// 30000 by default for the Minikube cluster
// 10000 by default for the Docker container
const PORT = __ENV.PORT || 10000;

export const options = {
  scenarios: {
    ramping_get_unsupported_path: {
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
    ramping_post_supported_path: {
      executor: "ramping-vus",
      startVus: 1,
      exec: "postPath",
      env: {
        path: "/v1/users",
        payload: JSON.stringify({
          username: "foo",
          password: "bar45678",
        }),
        params: JSON.stringify({
          headers: {
            "Content-Type": "application/json",
          },
        }),
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
  },
};

export function getPath() {
  const url = `http://${HOST}:${PORT}/${__ENV.path}`;
  http.get(url);
  sleep(1);
}

export function postPath() {
  const url = `http://${HOST}:${PORT}/${__ENV.path}`;
  const payload = __ENV.payload;
  const params = JSON.parse(__ENV.params);
  const res = http.post(url, payload, params);
  // TODO: Make this be scenario specific (say through __ENV defined in the scenario)
  check(res, {
    "is status 201 or 409": (r) => r.status === 200 || r.status === 409,
  });
  sleep(1);
}
