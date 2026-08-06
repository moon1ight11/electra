import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { loginRequest } from '../../api/auth';
import { usePhoneMask } from '../../hooks/usePhoneMask';
import styles from './Login.module.css';

export default function Login() {
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const phone = usePhoneMask();
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await loginRequest(phone.value, password);
      const meRes = await fetch('/api/v1/worker/me', { credentials: 'include' });
      const me = await meRes.json();
      login(me.id, me.name, me.role);
      navigate(me.role === 'owner' ? '/owner/dashboard' : '/dashboard');
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
            type="tel"
            value={phone.value}
            onChange={phone.onChange}
            placeholder="+7 (___) ___-__-__"
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