import styles from './Header.module.css';
import logo from '../../logo.svg';

export default function Header({ role, userName, onLogout, onMenuToggle }) {
  const roleLabel = role === 'owner' ? 'Владелец' : 'Исполнитель';

  return (
    <header className={styles.header}>
      <button className={styles.burger} onClick={onMenuToggle} aria-label="Меню">
        <span />
        <span />
        <span />
      </button>
      <img src={logo} alt="Electra" className={styles.logoImg} />
      <div className={styles.right}>
        <span className={styles.greeting}>{userName}</span>
        <span className={styles.role}>{roleLabel}</span>
        <button className={styles.logoutBtn} onClick={onLogout}>Выйти</button>
      </div>
    </header>
  );
}