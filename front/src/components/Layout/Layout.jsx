import { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import Header from './Header';
import Sidebar from './Sidebar';
import styles from './Layout.module.css';

export default function Layout({ children }) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <div className={styles.wrapper}>
      <Header
        role={user?.role}
        userName={user?.name}
        onLogout={handleLogout}
        onMenuToggle={() => setMenuOpen(!menuOpen)}
      />
      <div className={styles.body}>
        <Sidebar role={user?.role} open={menuOpen} onClose={() => setMenuOpen(false)} />
        <main className={styles.content} onClick={() => setMenuOpen(false)}>
          {children}
        </main>
      </div>
    </div>
  );
}