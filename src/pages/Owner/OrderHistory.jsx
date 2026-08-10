import { useState, useEffect } from 'react';
import OrderList from '../../components/OrderList/OrderList';
import styles from './OrderHistory.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

const BASE = process.env.REACT_APP_API_URL || '/api/v1';

export default function OrderHistory() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch(`${BASE}/owner/orders/history`, { credentials: 'include' });
        if (!res.ok) throw new Error('Ошибка загрузки');
        const data = await res.json();
        setOrders(data);
      } catch (err) {
        setError('Не удалось загрузить историю');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  return (
    <div className={styles.page}>
      <h1>История заказов</h1>
      <p className={styles.subtitle}>{orders.length} выполненных заказов</p>
      <OrderList orders={orders} emptyText="Нет выполненных заказов" showCompleted />
    </div>
  );
}