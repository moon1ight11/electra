const BASE = '/api/v1';

export async function fetchWorkers() {
  const res = await fetch(`${BASE}/worker/workers`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки исполнителей');
  return res.json();
}

export async function createWorker(name, phone, password, specialization) {
  const res = await fetch(`${BASE}/owner/workers`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ name, phone, password, specialization }),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Ошибка создания');
  }
  return res.json();
}