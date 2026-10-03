import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchWorkers } from '../../api/workers';
import { createDirectOrder } from '../../api/orders';
import styles from './CreateOrder.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

export default function CreateOrder() {
  const navigate = useNavigate();
  const [workers, setWorkers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [address, setAddress] = useState('');
  const [description, setDescription] = useState('');
  const [estimatedPrice, setEstimatedPrice] = useState('');
  const [plannedDate, setPlannedDate] = useState('');
  const [selectedWorkers, setSelectedWorkers] = useState([]);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        const data = await fetchWorkers();
        setWorkers(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  const toggleWorker = (id) => {
    setSelectedWorkers((prev) =>
      prev.includes(id) ? prev.filter((w) => w !== id) : [...prev, id]
    );
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (selectedWorkers.length === 0) {
      setError('Выберите хотя бы одного исполнителя');
      return;
    }
    setSaving(true);
    setError('');
    try {
      await createDirectOrder({
        address,
        description,
        estimated_price: estimatedPrice ? parseFloat(estimatedPrice) : null,
        planned_date: plannedDate || null,
        worker_ids: selectedWorkers,
      });
      navigate('/owner/orders/planned');
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <Skeleton count={4} />;

  return (
    <div className={styles.page}>
      <button className={styles.back} onClick={() => navigate(-1)}>← Назад</button>
      <h1>Новый заказ</h1>

      <form className={styles.form} onSubmit={handleSubmit}>
        {error && <p className={styles.error}>{error}</p>}

        <label>
          Адрес
          <input type="text" value={address} onChange={(e) => setAddress(e.target.value)} required />
        </label>
        <label>
          Описание
          <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={2} />
        </label>
        <div className={styles.formRow}>
          <label>
            Стоимость
            <input type="number" value={estimatedPrice} onChange={(e) => setEstimatedPrice(e.target.value)} />
          </label>
          <label>
            Дата
            <input type="date" value={plannedDate} onChange={(e) => setPlannedDate(e.target.value)} />
          </label>
        </div>

        <div className={styles.workerSection}>
          <span className={styles.workerLabel}>Исполнители</span>
          {workers.length === 0 ? (
            <span className={styles.noWorkers}>Нет исполнителей</span>
          ) : (
            <div className={styles.workerList}>
              {workers.map((w) => (
                <label key={w.id} className={styles.workerItem}>
                  <input
                    type="checkbox"
                    checked={selectedWorkers.includes(w.id)}
                    onChange={() => toggleWorker(w.id)}
                  />
                  <span>{w.name}</span>
                </label>
              ))}
            </div>
          )}
        </div>

        <button type="submit" className={styles.submitBtn} disabled={saving}>
          {saving ? 'Создание...' : 'Создать заказ'}
        </button>
      </form>
    </div>
  );
}