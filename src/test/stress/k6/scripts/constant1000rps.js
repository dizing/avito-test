import http from 'k6/http';
import { check } from 'k6';
import exec from 'k6/execution';

/*
  RPS — 1k SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%

  1000 req = 3 req/iteration * 333 iteration 

  VU cycle: -> login user, use all it's inventory on buying ->
*/
const BASE_URL = 'http://avito-shop-service:8080';
const globals = {
  VUcount: 40,
  PerVUIterations: 30000 / 40,
  BalanceRemains: 500 // reuse users in other stress tests
}

export const options = {
  scenarios: {
    buy: {
      exec: 'buy',
      executor: 'ramping-arrival-rate',
      startRate: 333,
      preAllocatedVUs: globals.VUcount,
      maxVUs: globals.VUcount,
      stages: [
        { target: 333, duration: '1m' }
      ]
    },
    sendCoins: {
      exec: 'send',
      executor: 'ramping-arrival-rate',
      startRate: 333,
      preAllocatedVUs: globals.VUcount,
      maxVUs: globals.VUcount,
      stages: [
        { target: 333, duration: '1m' }
      ]
    }
  },
  thresholds: {
    'http_req_duration{request:auth}': ['p(95)<100', 'p(99)<150'],
    'http_req_duration{request:buy}': ['p(99)<50'],
    'http_req_duration{request:info}': ['p(99)<50'],
    'http_req_duration{request:send}': ['p(99)<50'],
  }
};

function login(user_number) {
  const randomUsername = `vu_${user_number}_user`;
  const authPayload = JSON.stringify({
    username: randomUsername,
    password: 'pass'
  });

  const authParams = {
    headers: { 'Content-Type': 'application/json' },
    tags: { tag: 'auth' }
  };

  const authRes = http.post(`${BASE_URL}/api/auth`, authPayload, authParams);
  check(authRes, {
    'auth status is 200': (r) => r.status === 200,
  });

  return authRes.json().token;
}

function info(vuToken, expectedCoins, expectedSentLen) { // TODO: beatufy args. maybe several functions
  const params = {
    headers: {
      Authorization: `${vuToken}`
    },
    tags: { tag: 'info' }
  };

  let infoRes = http.get(`${BASE_URL}/api/info`, params);
  check(infoRes, {
    'info status is 200': (r) => r.status === 200,
  });

  if (expectedCoins !== undefined) {
    check(infoRes, {
      'user have expected coins count': (r) => r.json().coins === expectedCoins,
    });
  }

  if (expectedSentLen !== undefined) {
    check(infoRes, {
      'according transactions count': (r) => {
        if (!r.json().coinHistory.sent) {
          return expectedSentLen == 0;
        }

        return r.json().coinHistory.sent.len() === expectedSentLen;
      }
    });
  }
}

// ------ VU STATE BEGIN ------
var vuToken = null
var userNumberIterator = 0
var expectedBalance = 5000
var send_count = 0;

function getNumberForTestUser() {
  return __VU + exec.instance.vusInitialized * userNumberIterator
}

function ensureAuthorized() {
  if (vuToken === null) {
    if (__ITER >= globals.PerVUIterations) {
      throw new Error("Data exhausted!");
    }

    vuToken = login(getNumberForTestUser())
  }
}

function resetVuState() {
  vuToken = null;
  userNumberIterator = 0;
  expectedBalance = 5000;
  send_count = 0;
}

function nextVuIteration() {
  const iterator = userNumberIterator;
  resetVuState();
  userNumberIterator = iterator + 1;
}
// ------ VU STATE ENDS ------

export function buy() {
  const ITEM_NAME = "powerbank";
  const ITEM_PRICE = 200;

  if (__ITER === 0) {
    resetVuState()
  }

  ensureAuthorized();

  info(vuToken, expectedBalance, undefined)

  const params = {
    headers: {
      Authorization: `${vuToken}`
    },
    tags: { tag: 'buy' }
  };
  const buyRes = http.get(`${BASE_URL}/api/buy/${ITEM_NAME}`, params)
  check(buyRes, {
    'buy status is 200': (r) => r.status === ITEM_PRICE,
  });

  expectedBalance -= ITEM_PRICE;

  info(vuToken, expectedBalance, undefined)

  if (expectedBalance < ITEM_PRICE + globals.BalanceRemains) {
    nextVuIteration()
  }
}

export function send() {
  if (__ITER === 0) {
    resetVuState()
  }

  ensureAuthorized()

  const next_user_number = getNumberForTestUser() + 1;

  info(vuToken, undefined, send_count)

  const params = {
    headers: {
      Authorization: `${vuToken}`
    },
    tags: { tag: 'send' }
  };
  const sendPayload = JSON.stringify({
    toUser: `vu_${next_user_number}_user`,
    amount: 100
  });
  const sendRes = http.post(`${BASE_URL}/api/sendCoins`, sendPayload, params);
  check(sendRes, {
    'send status is 200': (r) => r.status === 200,
  });

  send_count++;

  info(vuToken, undefined, send_count)

  if (send_count == 5) {
    nextVuIteration()
  }
}