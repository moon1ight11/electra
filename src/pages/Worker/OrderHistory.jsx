import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import styles from './OrderHistory.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function OrderHistory() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadHistory();
  }, []);

  const loadHistory = async () => {
    try {
      const BASE = process.env.REACT_APP_API_URL || '/api/v1';
      const res = await fetch('${BASE}/public/requests`, {/worker/orders/history', { credentials: 'include' });
      if (!res.ok) throw new Error('Ошибка загрузки');
      const data = await res.json();
      setOrders(data);
    } catch (err) {
      setError('Не удалось загрузить историю');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  return (
    <div className={styles.page}>
      <h1>История заказов</h1>
      <p className={styles.subtitle}>{orders.length} выполненных заказов</p>

      {orders.length === 0 ? (
        <div className={styles.empty}>Пока нет выполненных заказов</div>
      ) : (
        <div className={styles.list}>
          {orders.map((order) => (
            <Link to={`/orders/${order.id}`} key={order.id} className={styles.card}>
              <div className={styles.cardLeft}>
                <span className={styles.address}>{order.address}</span>
                <p className={styles.description}>{order.description || 'Без описания'}</p>
              </div>
              <div className={styles.cardRight}>
                <span className={styles.date}>
                  {order.completed_at ? new Date(order.completed_at).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' }) : '—'}
                </span>
                {order.estimated_price && (
                  <span className={styles.price}>{order.estimated_price.toLocaleString()} ₽</span>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}