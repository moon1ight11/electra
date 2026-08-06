import { useState, useEffect, useCallback } from 'react';
import { createWorker, deleteWorker } from '../../api/workers';
import { usePhoneMask } from '../../hooks/usePhoneMask';
import Skeleton from '../../components/Skeleton/Skeleton';
import styles from './Workers.module.css';

const BASE = '/api/v1';

export default function Workers() {
  const [workers, setWorkers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [msg, setMsg] = useState('');

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const phone = usePhoneMask();
  const [password, setPassword] = useState('');
  const [specialization, setSpecialization] = useState('');

  const loadWorkers = useCallback(async () => {
    try {
      const [workersRes, statsRes] = await Promise.all([
        fetch(`${BASE}/worker/workers`, { credentials: 'include' }),
        fetch(`${BASE}/owner/statistics/all`, { credentials: 'include' }),
      ]);

      if (!workersRes.ok || !statsRes.ok) throw new Error('Ошибка загрузки');

      const workersData = await workersRes.json();
      const statsData = await statsRes.json();

      const merged = workersData.map((w) => {
        const stat = statsData.find((s) => s.worker_id === w.id);
        return {
          id: w.id,
          name: w.name,
          specialization: w.specialization || '',
          ordersCount: stat ? stat.orders_count : 0,
          totalEarned: stat ? stat.total_earned : 0,
          totalTimeSpent: stat ? stat.total_time_spent : 0,
        };
      });

      merged.sort((a, b) => b.totalEarned - a.totalEarned);
      setWorkers(merged);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadWorkers();
  }, [loadWorkers]);

  const handleCreate = async (e) => {
    e.preventDefault();
    setMsg('');
    try {
      await createWorker(name, phone.value, password, specialization);
      setMsg('Исполнитель создан');
      setShowForm(false);
      setName('');
      phone.reset();
      setPassword('');
      setSpecialization('');
      loadWorkers();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    }
  };

  const handleDelete = async (id, name) => {
    if (!window.confirm(`Удалить исполнителя "${name}"? Статистика сохранится.`)) return;
    setMsg('');
    try {
      await deleteWorker(id);
      setMsg('Исполнитель удалён');
      loadWorkers();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    }
  };

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  return (
    <div className={styles.page}>
      <div className={styles.headerRow}>
        <h1>Исполнители</h1>
        <button className={styles.addBtn} onClick={() => setShowForm(true)}>
          + Добавить
        </button>
      </div>

      {msg && (
        <p className={msg.includes('Ошибка') ? styles.error : styles.success}>{msg}</p>
      )}

      <div className={styles.list}>
        {workers.map((w, i) => (
          <div key={w.id} className={styles.card}>
            <div className={styles.rank}>
              <span className={i === 0 ? styles.rankGold : i === 1 ? styles.rankSilver : i === 2 ? styles.rankBronze : styles.rankDefault}>
                {i + 1}
              </span>
            </div>
            <div className={styles.cardBody}>
              <div className={styles.nameRow}>
                <span className={styles.workerName}>{w.name}</span>
                {w.specialization && (
                  <span className={styles.specialization}>{w.specialization}</span>
                )}
              </div>
              <div className={styles.stats}>
                <span>{w.ordersCount} заказов</span>
                <span>·</span>
                <span className={styles.earned}>{w.totalEarned.toLocaleString()} ₽</span>
                <span>·</span>
                <span>{Math.floor(w.totalTimeSpent / 60)} ч {w.totalTimeSpent % 60} мин</span>
              </div>
              <button
                className={styles.deleteBtn}
                onClick={() => handleDelete(w.id, w.name)}
              >
                Удалить
              </button>
            </div>
            {i === 0 && <span className={styles.topBadge}>🏆</span>}
          </div>
        ))}
      </div>

      {showForm && (
        <div className={styles.overlay} onClick={() => setShowForm(false)}>
          <form className={styles.form} onClick={(e) => e.stopPropagation()} onSubmit={handleCreate}>
            <h2>Новый исполнитель</h2>
            <label>
              Имя
              <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
            </label>
            <label>
              Телефон
              <input
                type="tel"
                value={phone.value}
                onChange={phone.onChange}
                placeholder="+7 (___) ___-__-__"
                required
              />
            </label>
            <label>
              Пароль
              <input type="text" value={password} onChange={(e) => setPassword(e.target.value)} required />
            </label>
            <label>
              Специализация
              <input type="text" value={specialization} onChange={(e) => setSpecialization(e.target.value)} placeholder="Электрик, отделочник..." />
            </label>
            <div className={styles.formActions}>
              <button type="submit" className={styles.submitBtn}>Создать</button>
              <button type="button" className={styles.cancelBtn} onClick={() => setShowForm(false)}>Отмена</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}