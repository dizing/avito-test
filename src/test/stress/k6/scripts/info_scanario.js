import http from 'k6/http';
import { check } from 'k6';

/*
  RPS — 1k SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%

  1000 req = 3 req/iteration * 333 iteration

  VU cycle: -> create user, use all it's inventory on buying ->
*/
export const options = {
  scenarios: {
    buy: {
      executor: 'per-vu-iterations',
      vus: 15,
      iterations: 20,
      maxDuration: '10s', // это не проверка, проверка достается из max http_req_duration
    },
  },
}

// per vu globals
var vuToken = null;
var expectedBalance = 5000

const BASE_URL = 'http://avito-shop-service:8080';

export default function () {
  if (vuToken === null) {
    const randomUsername = `vu_${__VU}_user_${Date.now()}`;
    const authPayload = JSON.stringify({
      username: randomUsername,
      password: 'pass'
    });

    const authParams = {
      headers: { 'Content-Type': 'application/json' },
    };

    const authRes = http.post(`${BASE_URL}/api/auth`, authPayload, authParams);
    vuToken = authRes.json().token;
  }

  const params = {
    headers: {
      Authorization: `${vuToken}`
    },
  };

  let infoRes = http.get(`${BASE_URL}/api/info`, params);
  check(infoRes, {
    'info status is 200': (r) => r.status === 200,
  });
  check(infoRes, {
    'base coins is 5000': (r) => r.json().coins === expectedBalance,
  });

  const buyRes = http.get(`${BASE_URL}/api/buy/powerbank`, params)
  check(buyRes, {
    'buy status is 200': (r) => r.status === 200,
  });
  
  expectedBalance -= 200;

  infoRes = http.get(`${BASE_URL}/api/info`, params);
  check(infoRes, {
    'info status is 200': (r) => r.status === 200,
  });
  check(infoRes, {
    'coins is updated after buying': (r) => r.json().coins === expectedBalance,
  });

  if (expectedBalance < 200) {
    vuToken = null;
    expectedBalance = 5000;
  }
}