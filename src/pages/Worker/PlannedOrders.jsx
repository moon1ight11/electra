import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { fetchPlannedOrders } from '../../api/orders';
import { fetchReports } from '../../api/reports';
import styles from './PlannedOrders.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function PlannedOrders() {
  const [orders, setOrders] = useState([]);
  const [reportStatuses, setReportStatuses] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadOrders();
  }, []);

  const loadOrders = async () => {
    try {
      const data = await fetchPlannedOrders();
      const sorted = data.sort((a, b) => new Date(a.planned_date) - new Date(b.planned_date));
      setOrders(sorted);

      const statuses = {};
      for (const order of sorted) {
        try {
          const reps = await fetchReports(order.id);
          const hasReport = reps.some((r) => r.time_spent || r.earned_amount);
          statuses[order.id] = hasReport;
        } catch {
          statuses[order.id] = false;
        }
      }
      setReportStatuses(statuses);
    } catch (err) {
      setError('Не удалось загрузить заказы');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  return (
    <div className={styles.page}>
      <h1>Заказы в работе</h1>
      <p className={styles.subtitle}>{orders.length} заказов</p>

      {orders.length === 0 ? (
        <div className={styles.empty}>Нет заказов в работе</div>
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
                  {order.planned_date ? new Date(order.planned_date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' }) : '—'}
                </span>
                {order.estimated_price && (
                  <span className={styles.price}>{order.estimated_price.toLocaleString()} ₽</span>
                )}
                {reportStatuses[order.id] ? (
                  <span className={styles.reportDone}>Отчёт заполнен</span>
                ) : (
                  <span className={styles.reportPending}>Заполните отчёт</span>
                )}
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}