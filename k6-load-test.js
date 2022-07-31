import http from "k6/http";
import { check, group, sleep } from "k6";
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
  stages: [
    { duration: "1m", target: 1000 }, // simulate ramp-up of traffic from 1 to 100 users over 5 minutes.
    { duration: "1m", target: 500 }, // stay at 100 users for 10 minutes
    { duration: "1m", target: 0 }, // ramp-down to 0 users
  ],
  thresholds: {
    http_req_duration: ["p(99)<1500"], // 99% of requests must complete below 1.5s
  },
};

const BASE_URL = "http://localhost";

export default () => {
  const response = http.get(`${BASE_URL}/server-info`);
  check(response, { "status was 200": (r) => r.status == 200 });
  sleep(0.5);
};

export function handleSummary(data) {
  return {
    "summary.html": htmlReport(data),
  };
}
