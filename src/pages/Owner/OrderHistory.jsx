import { useState, useEffect } from 'react';
import OrderList from '../../components/OrderList/OrderList';
import styles from './OrderHistory.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function OrderHistory() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch('/api/v1/owner/orders/history', { credentials: 'include' });
        if (!res.ok) throw new Error('Ошибка');
        const data = await res.json();
        setOrders(data);
      } catch { } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <Skeleton count={4} />;

  return (
    <div className={styles.page}>
      <h1>История заказов</h1>
      <p className={styles.subtitle}>{orders.length} выполненных заказов</p>
      <OrderList orders={orders} emptyText="Нет выполненных заказов" showCompleted />
    </div>
  );
}