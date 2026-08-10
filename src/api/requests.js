const BASE = process.env.REACT_APP_API_URL || '/api/v1';

export async function fetchNewRequests() {
  const res = await fetch(`${BASE}/owner/requests/new`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки заявок');
  return res.json();
}

export async function fetchAllRequests() {
  const res = await fetch(`${BASE}/owner/requests/all`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки заявок');
  return res.json();
}

export async function cancelRequest(id) {
  const res = await fetch(`${BASE}/owner/requests/${id}/cancel`, {
    method: 'PATCH',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка отмены заявки');
  return res.json();
}

export async function convertRequest(data) {
  const res = await fetch(`${BASE}/owner/requests/convert`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Ошибка создания заказа');
  }
  return res.json();
}