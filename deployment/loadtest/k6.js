import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 2000,
  duration: '5m',
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL = 'http://localhost:3000'; // api-gateway

export default function () {
  // Health check
  let res = http.get(`${BASE_URL}/health`);
  check(res, { 'health 200': (r) => r.status === 200 });

  // Login simulation (adjust payload)
  let loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
    email: 'test@example.com',
    password: 'testpass'
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(loginRes, { 'login ok': (r) => r.status === 200 });

  // Refresh token simulation
  if (loginRes.json('token')) {
    let refreshRes = http.post(`${BASE_URL}/auth/refresh`, {}, {
      headers: { 'Authorization': `Bearer ${loginRes.json('token')}` },
    });
    check(refreshRes, { 'refresh ok': (r) => r.status === 200 });
  }

  // Job list
  http.get(`${BASE_URL}/gov-jobs`);
  http.get(`${BASE_URL}/private-jobs`);

  sleep(1);
}
