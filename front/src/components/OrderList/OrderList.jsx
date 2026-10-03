import { Link } from 'react-router-dom';
import styles from './OrderList.module.css';

export default function OrderList({ orders, emptyText = 'Нет заказов', showCompleted = false }) {
  if (orders.length === 0) {
    return <div className={styles.empty}>{emptyText}</div>;
  }

  return (
    <div className={styles.list}>
      {orders.map((order) => (
        <Link to={`/owner/orders/${order.id}`} key={order.id} className={styles.card}>
          <div className={styles.cardLeft}>
            <span className={styles.address}>{order.address}</span>
            <p className={styles.description}>{order.description || 'Без описания'}</p>
          </div>
          <div className={styles.cardRight}>
            <span className={styles.date}>
              {showCompleted && order.completed_at
                ? new Date(order.completed_at).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })
                : order.planned_date
                  ? new Date(order.planned_date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })
                  : '—'}
            </span>
            {order.estimated_price && (
              <span className={styles.price}>{order.estimated_price.toLocaleString()} ₽</span>
            )}
            {showCompleted && (
              <span className={styles.badge}>Завершён</span>
            )}
          </div>
        </Link>
      ))}
    </div>
  );
}