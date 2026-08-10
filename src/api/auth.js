const BASE = process.env.REACT_APP_API_URL || '/api/v1';

export async function loginRequest(phone, password) {
  const res = await fetch(`${BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ phone, password }),
    credentials: 'include',
  });
  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.error || 'Ошибка входа');
  }
  return res.json();
}

export async function logoutRequest() {
  await fetch(`${BASE}/worker/logout`, {
    method: 'POST',
    credentials: 'include',
  });
}