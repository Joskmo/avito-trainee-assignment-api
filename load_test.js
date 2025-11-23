import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  stages: [
    { duration: '10s', target: 10 }, // Разгон
    { duration: '30s', target: 50 }, // Нагрузка 50 пользователей
    { duration: '10s', target: 0 },  // Спад
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'], // 95% запросов должны быть быстрее 100мс
    http_req_failed: ['rate<0.01'],   // Ошибок меньше 1%
  },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  const uniqueId = randomString(8);
  const teamName = `load_team_${uniqueId}`;
  const user1 = `u1_${uniqueId}`;
  const user2 = `u2_${uniqueId}`;
  const user3 = `u3_${uniqueId}`;

  // 1. Создание команды
  const createTeamPayload = JSON.stringify({
    team_name: teamName,
    members: [
      { user_id: user1, username: "User1", is_active: true },
      { user_id: user2, username: "User2", is_active: true },
      { user_id: user3, username: "User3", is_active: true },
    ],
  });

  let res = http.post(`${BASE_URL}/team/add`, createTeamPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, { 'team created': (r) => r.status === 201 });

  // 2. Создание PR
  const createPRPayload = JSON.stringify({
    pull_request_id: `pr_${uniqueId}`,
    pull_request_name: `PR ${uniqueId}`,
    author_id: user1,
  });

  res = http.post(`${BASE_URL}/pullRequest/create`, createPRPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  
  const prCreated = check(res, { 'pr created': (r) => r.status === 201 });

  if (prCreated) {
    // 3. Массовая деактивация (самая тяжелая операция)
    const deactivatePayload = JSON.stringify({
        users: [user2]
    });
    
    res = http.post(`${BASE_URL}/team/deactivateUsers`, deactivatePayload, {
        headers: { 'Content-Type': 'application/json' },
    });
    check(res, { 'deactivate success': (r) => r.status === 200 });
  }

  // 4. Статистика
  res = http.get(`${BASE_URL}/stats`);
  check(res, { 'stats ok': (r) => r.status === 200 });

  sleep(1);
}
