const BASE = '/api/v1';

export async function fetchAllStats(from = '', to = '') {
  const params = new URLSearchParams();
  if (from) params.set('from', from);
  if (to) params.set('to', to);

  const res = await fetch(`${BASE}/owner/statistics/all?${params}`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки статистики');
  return res.json();
}

export async function fetchWorkerStats(workerId, from = '', to = '') {
  const params = new URLSearchParams();
  if (from) params.set('from', from);
  if (to) params.set('to', to);

  const res = await fetch(`${BASE}/owner/statistics/workers/${workerId}?${params}`, {
    credentials: 'include',
  });
  if (!res.ok) throw new Error('Ошибка загрузки статистики');
  return res.json();
}