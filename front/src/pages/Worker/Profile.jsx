import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { updateProfile } from '../../api/workers';
import styles from './Profile.module.css';

export default function Profile() {
  const { user, login } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState(user?.name || '');
  const [specialization, setSpecialization] = useState('');
  const [password, setPassword] = useState('');
  const [msg, setMsg] = useState('');
  const [saving, setSaving] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    setMsg('');
    try {
      const updated = await updateProfile({
        name,
        specialization,
        password: password || undefined,
      });
      login(updated.id, updated.name, updated.role);
      setMsg('Профиль обновлён');
      setPassword('');
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className={styles.page}>
      <button className={styles.back} onClick={() => navigate(-1)}>← Назад</button>
      <h1>Редактирование профиля</h1>

      <form className={styles.form} onSubmit={handleSubmit}>
        {msg && (
          <p className={msg.includes('Ошибка') ? styles.error : styles.success}>{msg}</p>
        )}

        <label>
          Имя
          <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
        </label>

        <label>
          Специализация
          <input
            type="text"
            value={specialization}
            onChange={(e) => setSpecialization(e.target.value)}
            placeholder="Электрик, отделочник..."
          />
        </label>

        <label>
          Новый пароль (оставьте пустым, если не меняете)
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Новый пароль"
          />
        </label>

        <button type="submit" disabled={saving} className={styles.saveBtn}>
          {saving ? 'Сохранение...' : 'Сохранить'}
        </button>
      </form>
    </div>
  );
}