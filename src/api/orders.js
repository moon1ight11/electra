const BASE = '/api/v1';

export async function fetchPlannedOrders() {
  const res = await fetch(`${BASE}/worker/orders/planned`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки заказов');
  return res.json();
}

export async function fetchAllPlannedOrders() {
  const res = await fetch(`${BASE}/owner/orders/planned`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки заказов');
  return res.json();
}

export async function createDirectOrder(data) {
  const res = await fetch(`${BASE}/owner/orders/direct`, {
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