import { useState, useEffect } from 'react';
import { fetchAllPlannedOrders } from '../../api/orders';
import OrderList from '../../components/OrderList/OrderList';
import styles from './PlannedOrders.module.css';
import { Link } from 'react-router-dom';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function PlannedOrders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      try {
        const data = await fetchAllPlannedOrders();
        setOrders(data.sort((a, b) => new Date(a.planned_date) - new Date(b.planned_date)));
      } catch { } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) return <Skeleton count={4} />;

  return (
    <div className={styles.page}>
      <h1>Запланированные заказы</h1>
      <p className={styles.subtitle}>{orders.length} заказов</p>
      <Link to="/owner/orders/create" className={styles.createBtn}>+ Новый заказ</Link>
      <OrderList orders={orders} emptyText="Нет запланированных заказов" />
    </div>
  );
}