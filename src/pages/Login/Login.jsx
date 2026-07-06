import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { loginRequest } from '../../api/auth';
import styles from './Login.module.css';

export default function Login() {
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await loginRequest(phone, password);

      const [ownerCheck, meRes] = await Promise.all([
        fetch('/api/v1/owner/requests/new'),
        fetch('/api/v1/worker/me', { credentials: 'include' }),
      ]);

      const role = ownerCheck.ok ? 'owner' : 'worker';
      const me = await meRes.json();

      login(me.id, me.name, role);
      navigate(role === 'owner' ? '/owner/dashboard' : '/dashboard');
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className={styles.container}>
      <Link to="/" className={styles.backLink}>← На главную</Link>
      <form className={styles.form} onSubmit={handleSubmit}>
        <h1>Вход</h1>
        {error && <p className={styles.error}>{error}</p>}
        <label>
          Телефон
          <input
            type="text"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            placeholder="79991234567"
            required
          />
        </label>
        <label>
          Пароль
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </label>
        <button type="submit">Войти</button>
      </form>
    </div>
  );
}