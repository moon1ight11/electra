const BASE = '/api/v1';

export async function fetchReports(orderId) {
  const res = await fetch(`${BASE}/worker/orders/${orderId}/reports`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки отчётов');
  return res.json();
}

export async function updateReport(orderId, data) {
  const res = await fetch(`${BASE}/worker/orders/report`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ order_id: orderId, ...data }),
  });
  if (!res.ok) throw new Error('Ошибка обновления отчёта');
  return res.json();
}

export async function completeOrder(orderId) {
  const res = await fetch(`${BASE}/worker/orders/${orderId}/complete`, {
    method: 'PATCH',
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка завершения заказа');
  return res.json();
}