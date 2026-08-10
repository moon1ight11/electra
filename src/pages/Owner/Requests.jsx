import { useState, useEffect, useCallback } from 'react';
import { fetchNewRequests, fetchAllRequests, cancelRequest, convertRequest } from '../../api/requests';
import styles from './Requests.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';


export default function Requests() {
  const [tab, setTab] = useState('new');
  const [requests, setRequests] = useState([]);
  const [workers, setWorkers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [msg, setMsg] = useState('');

  const [showForm, setShowForm] = useState(null);
  const [address, setAddress] = useState('');
  const [description, setDescription] = useState('');
  const [estimatedPrice, setEstimatedPrice] = useState('');
  const [plannedDate, setPlannedDate] = useState('');
  const [selectedWorkers, setSelectedWorkers] = useState([]);

  const loadData = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const fetcher = tab === 'new' ? fetchNewRequests : fetchAllRequests;
      const [reqData] = await Promise.all([fetcher()]);
      setRequests(reqData);

      const BASE = process.env.REACT_APP_API_URL || '/api/v1';
      const res = await fetch('${BASE}worker/workers', { credentials: 'include' });
      if (res.ok) {
        const workers = await res.json();
        setWorkers(workers);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [tab]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleCancel = async (id) => {
    if (!window.confirm('Отменить заявку?')) return;
    try {
      await cancelRequest(id);
      setMsg('Заявка отменена');
      loadData();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    }
  };

  const openForm = (req) => {
    setShowForm(req.id);
    setAddress('');
    setDescription(req.comment || '');
    setEstimatedPrice('');
    setPlannedDate('');
    setSelectedWorkers([]);
  };

  const toggleWorker = (id) => {
    setSelectedWorkers((prev) =>
      prev.includes(id) ? prev.filter((w) => w !== id) : [...prev, id]
    );
  };

  const handleConvert = async (e) => {
    e.preventDefault();
    if (selectedWorkers.length === 0) {
      setMsg('Выберите хотя бы одного исполнителя');
      return;
    }
    setMsg('');
    try {
      await convertRequest({
        request_id: showForm,
        address,
        description,
        estimated_price: estimatedPrice ? parseFloat(estimatedPrice) : null,
        planned_date: plannedDate || null,
        worker_ids: selectedWorkers,
      });
      setMsg('Заказ создан');
      setShowForm(null);
      loadData();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    }
  };

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;

  return (
    <div className={styles.page}>
      <h1>Заявки</h1>

      <div className={styles.tabs}>
        <button className={tab === 'new' ? styles.tabActive : styles.tab} onClick={() => setTab('new')}>Новые</button>
        <button className={tab === 'all' ? styles.tabActive : styles.tab} onClick={() => setTab('all')}>Все</button>
      </div>

      {msg && <p className={msg.includes('Ошибка') ? styles.error : styles.success}>{msg}</p>}

      {requests.length === 0 ? (
        <div className={styles.empty}>Нет заявок</div>
      ) : (
        <div className={styles.list}>
          {requests.map((r) => (
            <div key={r.id} className={styles.card}>
              <div className={styles.cardBody}>
                <div className={styles.cardHeader}>
                  <span className={styles.name}>{r.name}</span>
                  <span className={styles.phone}>{r.phone}</span>
                  {r.status !== 'new' && (
                    <span className={styles.badge}>{r.status === 'converted' ? 'В заказе' : 'Отменена'}</span>
                  )}
                </div>
                <p className={styles.comment}>{r.comment || 'Без комментария'}</p>
                <span className={styles.date}>
                  {new Date(r.created_at).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', hour: '2-digit', minute: '2-digit' })}
                </span>
              </div>
              {r.status === 'new' && (
                <div className={styles.actions}>
                  <button className={styles.convertBtn} onClick={() => openForm(r)}>Создать заказ</button>
                  <button className={styles.cancelBtn} onClick={() => handleCancel(r.id)}>Отменить</button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <div className={styles.overlay} onClick={() => setShowForm(null)}>
          <form className={styles.form} onClick={(e) => e.stopPropagation()} onSubmit={handleConvert}>
            <h2>Создать заказ</h2>
            <label>
              Адрес
              <input type="text" value={address} onChange={(e) => setAddress(e.target.value)} required />
            </label>
            <label>
              Описание
              <input type="text" value={description} onChange={(e) => setDescription(e.target.value)} />
            </label>
            <label>
              Стоимость
              <input type="number" value={estimatedPrice} onChange={(e) => setEstimatedPrice(e.target.value)} />
            </label>
            <label>
              Дата
              <input type="date" value={plannedDate} onChange={(e) => setPlannedDate(e.target.value)} />
            </label>

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

            <div className={styles.formActions}>
              <button type="submit" className={styles.convertBtn}>Создать</button>
              <button type="button" className={styles.cancelBtn} onClick={() => setShowForm(null)}>Отмена</button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}