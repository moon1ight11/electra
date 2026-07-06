import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { fetchAllPlannedOrders } from '../../api/orders';
import { fetchReports, updateReport } from '../../api/reports';
import styles from './OrderDetail.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

const BASE = '/api/v1';

export default function OrderDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [order, setOrder] = useState(null);
  const [reports, setReports] = useState([]);
  const [workerNames, setWorkerNames] = useState({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [msg, setMsg] = useState('');
  const [removing, setRemoving] = useState(false);
  const [completing, setCompleting] = useState(false);

  const [timeSpent, setTimeSpent] = useState('');
  const [earnedAmount, setEarnedAmount] = useState('');
  const [materials, setMaterials] = useState('');
  const [notes, setNotes] = useState('');
  const [saving, setSaving] = useState(false);

  const loadWorkerNames = useCallback(async () => {
    try {
      const res = await fetch('/api/v1/worker/workers', { credentials: 'include' });
      if (res.ok) {
        const workers = await res.json();
        const names = {};
        workers.forEach((w) => {
          names[w.id] = w.name;
        });
        setWorkerNames(names);
      }
    } catch { }
  }, []);

  const loadOrder = useCallback(async () => {
    try {
      const planned = await fetchAllPlannedOrders();
      let found = planned.find((o) => o.id === id);
      if (!found) {
        const histRes = await fetch(`${BASE}/owner/orders/history`, { credentials: 'include' });
        const hist = await histRes.json();
        found = hist.find((o) => o.id === id);
        if (!found) throw new Error('Заказ не найден');
      }
      setOrder(found);
    } catch (err) {
      setError(err.message);
    }
  }, [id]);

  const loadReports = useCallback(async () => {
    try {
      const reps = await fetchReports(id);
      setReports(reps);
      if (user) {
        const myReport = reps.find((r) => r.worker_id === user.id);
        if (myReport) {
          if (myReport.time_spent) setTimeSpent(String(myReport.time_spent));
          if (myReport.earned_amount) setEarnedAmount(String(myReport.earned_amount));
          if (myReport.materials_used) setMaterials(myReport.materials_used);
          if (myReport.notes) setNotes(myReport.notes);
        }
      }
    } catch { }
  }, [id, user]);

  useEffect(() => {
    const init = async () => {
      setLoading(true);
      setError('');
      await Promise.all([loadOrder(), loadWorkerNames()]);
      await loadReports();
      setLoading(false);
    };
    init();
  }, [loadOrder, loadReports, loadWorkerNames]);

  const handleSaveReport = async (e) => {
    e.preventDefault();
    setSaving(true);
    setMsg('');
    try {
      if (!user) {
        setMsg('Ошибка: пользователь не определён');
        return;
      }
      const myReport = reports.find((r) => r.worker_id === user.id);
      if (!myReport) {
        setMsg('Ошибка: вы не назначены на этот заказ');
        return;
      }
      await updateReport(id, {
        time_spent: timeSpent ? Number(timeSpent) : null,
        earned_amount: earnedAmount ? parseFloat(earnedAmount) : null,
        materials_used: materials || null,
        notes: notes || null,
      });
      setMsg('Отчёт сохранён');
      await loadReports();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    } finally {
      setSaving(false);
    }
  };

  const handleRemoveWorker = async (workerId) => {
    if (!window.confirm('Снять исполнителя с заказа?')) return;
    setRemoving(true);
    setMsg('');
    try {
      const res = await fetch(`${BASE}/owner/orders/${id}/workers/${workerId}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      if (!res.ok) throw new Error('Ошибка');
      setMsg('Исполнитель снят');
      await loadReports();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    } finally {
      setRemoving(false);
    }
  };

  const handleComplete = async () => {
    if (!window.confirm('Завершить заказ?')) return;
    setCompleting(true);
    setMsg('');
    try {
      const res = await fetch(`${BASE}/owner/orders/${id}/complete`, {
        method: 'PATCH',
        credentials: 'include',
      });
      if (!res.ok) throw new Error('Ошибка');
      navigate('/owner/orders/history');
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
      setCompleting(false);
    }
  };

  if (loading) return <Skeleton count={4} />;
  if (error) return <div className={styles.status}>{error}</div>;
  if (!order) return <div className={styles.status}>Заказ не найден</div>;

  const isCompleted = !!order.completed_at;

  return (
    <div className={styles.page}>
      <button className={styles.back} onClick={() => navigate(-1)}>← Назад</button>

      <div className={styles.header}>
        <h1>{order.address}</h1>
        {isCompleted && <span className={styles.badge}>Завершён</span>}
      </div>

      <div className={styles.infoGrid}>
        <div className={styles.infoItem}>
          <span className={styles.label}>Описание</span>
          <span>{order.description || '—'}</span>
        </div>
        <div className={styles.infoItem}>
          <span className={styles.label}>Дата</span>
          <span>{order.planned_date ? new Date(order.planned_date).toLocaleDateString('ru-RU') : '—'}</span>
        </div>
        <div className={styles.infoItem}>
          <span className={styles.label}>Стоимость</span>
          <span className={styles.price}>{order.estimated_price ? order.estimated_price.toLocaleString() + ' ₽' : '—'}</span>
        </div>
        {order.request_id && (
          <div className={styles.infoItem}>
            <span className={styles.label}>Заявка</span>
            <span className={styles.requestId}>{order.request_id}</span>
          </div>
        )}
      </div>

      <h2 className={styles.sectionTitle}>Исполнители</h2>
      {reports.length === 0 ? (
        <p className={styles.empty}>Нет назначенных исполнителей</p>
      ) : (
        <div className={styles.reports}>
          {reports.map((r) => (
            <div key={r.id} className={styles.reportCard}>
              <div className={styles.reportHead}>
                <span className={styles.workerName}>{workerNames[r.worker_id] || 'Неизвестный'}</span>
                {!isCompleted && (
                  <button
                    className={styles.removeBtn}
                    onClick={() => handleRemoveWorker(r.worker_id)}
                    disabled={removing}
                  >
                    Снять
                  </button>
                )}
              </div>
              <div className={styles.reportRow}>
                <span className={styles.reportLabel}>Время</span>
                <span>{r.time_spent ? r.time_spent + ' мин' : '—'}</span>
              </div>
              <div className={styles.reportRow}>
                <span className={styles.reportLabel}>Заработок</span>
                <span className={styles.price}>{r.earned_amount ? r.earned_amount.toLocaleString() + ' ₽' : '—'}</span>
              </div>
              <div className={styles.reportRow}>
                <span className={styles.reportLabel}>Материалы</span>
                <span>{r.materials_used || '—'}</span>
              </div>
              <div className={styles.reportRow}>
                <span className={styles.reportLabel}>Заметки</span>
                <span>{r.notes || '—'}</span>
              </div>
            </div>
          ))}
        </div>
      )}

      {!isCompleted && (
        <>
          <h2 className={styles.sectionTitle}>Мой отчёт</h2>
          <form className={styles.form} onSubmit={handleSaveReport}>
            <div className={styles.formRow}>
              <label>
                Время (мин)
                <input type="number" value={timeSpent} onChange={(e) => setTimeSpent(e.target.value)} />
              </label>
              <label>
                Заработок (₽)
                <input type="number" value={earnedAmount} onChange={(e) => setEarnedAmount(e.target.value)} />
              </label>
            </div>
            <label>
              Материалы
              <input type="text" value={materials} onChange={(e) => setMaterials(e.target.value)} placeholder="Кабель, автоматы..." />
            </label>
            <label>
              Заметки
              <textarea value={notes} onChange={(e) => setNotes(e.target.value)} rows={2} placeholder="Что сделано, особенности..." />
            </label>
            <button type="submit" disabled={saving} className={styles.saveBtn}>
              {saving ? 'Сохранение...' : 'Сохранить отчёт'}
            </button>
            {msg && <p className={msg.includes('Ошибка') ? styles.error : styles.success}>{msg}</p>}
          </form>
        </>
      )}

      {!isCompleted && (
        <button
          className={styles.completeBtn}
          onClick={handleComplete}
          disabled={completing}
        >
          {completing ? 'Завершение...' : 'Завершить заказ'}
        </button>
      )}
    </div>
  );
}