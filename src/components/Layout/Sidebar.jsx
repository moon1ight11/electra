import { NavLink } from 'react-router-dom';
import styles from './Sidebar.module.css';

const workerLinks = [
  { to: '/dashboard', label: 'Главная' },
  { to: '/orders/planned', label: 'В работе' },
  { to: '/orders/history', label: 'История' },
  { to: '/help', label: 'Помощь' },
  { to: '/profile', label: 'Профиль' },
];

const ownerLinks = [
  { to: '/owner/dashboard', label: 'Главная' },
  { to: '/owner/requests', label: 'Заявки' },
  { to: '/owner/orders/planned', label: 'Заказы' },
  { to: '/owner/orders/history', label: 'История' },
  { to: '/owner/statistics', label: 'Статистика' },
  { to: '/owner/workers', label: 'Исполнители' },
];

export default function Sidebar({ role, open, onClose }) {
  const links = role === 'owner' ? ownerLinks : workerLinks;

  return (
    <>
      {/* Оверлей */}
      {open && <div className={styles.overlay} onClick={onClose} />}

      <aside className={`${styles.sidebar} ${open ? styles.open : ''}`}>
        <nav className={styles.nav}>
          {links.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              className={({ isActive }) => isActive ? styles.active : styles.link}
              onClick={onClose}
            >
              {link.label}
            </NavLink>
          ))}
        </nav>
      </aside>
    </>
  );
}