import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { fetchNewRequests } from '../../api/requests';
import { fetchAllPlannedOrders } from '../../api/orders';
import styles from './Dashboard.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function Dashboard() {
  const [requests, setRequests] = useState([]);
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const load = async () => {
      try {
        const [reqData, ordData] = await Promise.all([
          fetchNewRequests(),
          fetchAllPlannedOrders(),
        ]);
        setRequests(reqData);
        setOrders(ordData);
      } catch (err) {
        setError('Не удалось загрузить данные');
      }
      setLoading(false);
    };
    load();
  }, []);

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  const newRequests = requests.slice(0, 5);
  const upcomingOrders = orders.slice(0, 5);

  return (
    <div className={styles.dashboard}>
      <h1>Дашборд</h1>

      <div className={styles.cards}>
        <div className={styles.statCard}>
          <span className={styles.statNumber}>{requests.length}</span>
          <span className={styles.statLabel}>Новых заявок</span>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statNumber}>{orders.length}</span>
          <span className={styles.statLabel}>Заказов в работе</span>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statNumber}>
            {orders.reduce((sum, o) => sum + (o.estimated_price || 0), 0).toLocaleString()} ₽
          </span>
          <span className={styles.statLabel}>Общая сумма</span>
        </div>
      </div>

      <div className={styles.columns}>
        <div className={styles.column}>
          <div className={styles.columnHeader}>
            <h2>Новые заявки</h2>
            <Link to="/owner/requests" className={styles.allLink}>Все →</Link>
          </div>
          {newRequests.length === 0 ? (
            <p className={styles.empty}>Нет новых заявок</p>
          ) : (
            newRequests.map((r) => (
              <div key={r.id} className={styles.miniCard}>
                <div>
                  <span className={styles.name}>{r.name}</span>
                  <span className={styles.phone}>{r.phone}</span>
                </div>
                <p className={styles.comment}>{r.comment || '—'}</p>
              </div>
            ))
          )}
        </div>

        <div className={styles.column}>
          <div className={styles.columnHeader}>
            <h2>Ближайшие заказы</h2>
            <Link to="/owner/orders/planned" className={styles.allLink}>Все →</Link>
          </div>
          {upcomingOrders.length === 0 ? (
            <p className={styles.empty}>Нет запланированных заказов</p>
          ) : (
            upcomingOrders.map((o) => (
              <Link to={`/owner/orders/${o.id}`} key={o.id} className={styles.miniCard}>
                <div className={styles.orderRow}>
                  <span className={styles.address}>{o.address}</span>
                  <span className={styles.date}>
                    {o.planned_date ? new Date(o.planned_date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' }) : '—'}
                  </span>
                </div>
                {o.estimated_price && (
                  <span className={styles.price}>{o.estimated_price.toLocaleString()} ₽</span>
                )}
              </Link>
            ))
          )}
        </div>
      </div>
    </div>
  );
}