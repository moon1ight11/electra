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

export async function deleteWorker(id) {
  const res = await fetch(`${BASE}/owner/workers/${id}`, {
    method: 'DELETE',
    credentials: 'include',
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Ошибка удаления');
  }
  return res.json();
}

export async function updateProfile(data) {
  const res = await fetch(`${BASE}/worker/me`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || 'Ошибка обновления профиля');
  }
  return res.json();
}