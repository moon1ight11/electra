import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { fetchPlannedOrders } from '../../api/orders';
import { fetchReports } from '../../api/reports';
import styles from './Dashboard.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function Dashboard() {
  const { user } = useAuth();
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
      setOrders(data);

      const statuses = {};
      for (const order of data) {
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

  const upcoming = orders.slice(0, 5);

  return (
    <div className={styles.dashboard}>
      <h1>Добрый день, {user?.name}!</h1>
      <p className={styles.subtitle}>У вас {orders.length} заказов в работе</p>

      {upcoming.length === 0 ? (
        <div className={styles.empty}>
          <p>Пока нет заказов в работе</p>
        </div>
      ) : (
        <div className={styles.grid}>
          {upcoming.map((order) => (
            <Link to={`/orders/${order.id}`} key={order.id} className={styles.card}>
              <div className={styles.cardHeader}>
                <span className={styles.address}>{order.address}</span>
                <span className={styles.date}>
                  {order.planned_date ? new Date(order.planned_date).toLocaleDateString('ru-RU') : '—'}
                </span>
              </div>
              <p className={styles.description}>{order.description || 'Без описания'}</p>
              <div className={styles.cardFooter}>
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

      {orders.length > 5 && (
        <Link to="/orders/planned" className={styles.allLink}>
          Все заказы →
        </Link>
      )}
    </div>
  );
}