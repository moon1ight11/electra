import { useState, useEffect } from 'react';
import styles from './Statistics.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function Statistics() {
  const [stats, setStats] = useState([]);
  const [summary, setSummary] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async (fromDate = '', toDate = '') => {
    setLoading(true);
    setError('');
    try {
      const params = new URLSearchParams();
      if (fromDate) params.set('from', fromDate);
      if (toDate) params.set('to', toDate);

      const [statsRes, summaryRes] = await Promise.all([
        fetch(`/api/v1/owner/statistics/all?${params}`, { credentials: 'include' }),
        fetch(`/api/v1/owner/statistics/summary?${params}`, { credentials: 'include' }),
      ]);

      if (!statsRes.ok || !summaryRes.ok) throw new Error('Ошибка загрузки');

      const statsData = await statsRes.json();
      const summaryData = await summaryRes.json();

      setStats(statsData);
      setSummary(summaryData);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleFilter = (e) => {
    e.preventDefault();
    loadStats(from, to);
  };

  const handleReset = () => {
    setFrom('');
    setTo('');
    loadStats();
  };

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  return (
    <div className={styles.page}>
      <h1>Статистика</h1>

      <form className={styles.filter} onSubmit={handleFilter}>
        <label>
          С
          <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
        </label>
        <label>
          По
          <input type="date" value={to} onChange={(e) => setTo(e.target.value)} />
        </label>
        <button type="submit" className={styles.filterBtn}>Применить</button>
        <button type="button" className={styles.resetBtn} onClick={handleReset}>Сброс</button>
      </form>

      {summary && (
        <div className={styles.summary}>
          <div className={styles.summaryCard}>
            <span className={styles.summaryNumber}>{summary.orders_count}</span>
            <span className={styles.summaryLabel}>Заказов</span>
          </div>
          <div className={styles.summaryCard}>
            <span className={styles.summaryNumber}>{summary.total_earned.toLocaleString()} ₽</span>
            <span className={styles.summaryLabel}>Заработано</span>
          </div>
          <div className={styles.summaryCard}>
            <span className={styles.summaryNumber}>{Math.floor(summary.total_time_spent / 60)} ч</span>
            <span className={styles.summaryLabel}>Потрачено времени</span>
          </div>
        </div>
      )}

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Исполнитель</th>
              <th>Заказов</th>
              <th>Заработано</th>
              <th>Время</th>
            </tr>
          </thead>
          <tbody>
            {stats.length === 0 ? (
              <tr>
                <td colSpan={4} className={styles.emptyCell}>Нет данных</td>
              </tr>
            ) : (
              stats.map((s) => (
                <tr key={s.worker_id}>
                  <td className={styles.nameCell} data-label="Исполнитель">{s.worker_name}</td>
                  <td data-label="Заказов">{s.orders_count}</td>
                  <td className={styles.priceCell} data-label="Заработано">{s.total_earned.toLocaleString()} ₽</td>
                  <td data-label="Время">{Math.floor(s.total_time_spent / 60)} ч {s.total_time_spent % 60} мин</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}