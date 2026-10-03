import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { fetchAllPlannedOrders } from '../../api/orders';
import { fetchReports, updateReport } from '../../api/reports';
import styles from './OrderDetail.module.css';
import Skeleton from '../../components/Skeleton/Skeleton';

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

  const [editMode, setEditMode] = useState(false);
  const [editAddress, setEditAddress] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editEstimatedPrice, setEditEstimatedPrice] = useState('');
  const [editPlannedDate, setEditPlannedDate] = useState('');
  const [editSaving, setEditSaving] = useState(false);

  const BASE = process.env.REACT_APP_API_URL || '/api/v1';

  const loadWorkerNames = useCallback(async () => {
    try {
      const res = await fetch(`${BASE}/worker/workers`, { credentials: 'include' });
      if (res.ok) {
        const workers = await res.json();
        const names = {};
        workers.forEach((w) => {
          names[w.id] = w.name;
        });
        setWorkerNames(names);
      }
    } catch { }
  }, [BASE]);

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
  }, [id, BASE]);

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

  const openEditMode = () => {
    setEditAddress(order.address || '');
    setEditDescription(order.description || '');
    setEditEstimatedPrice(order.estimated_price ? String(order.estimated_price) : '');
    setEditPlannedDate(order.planned_date ? order.planned_date.slice(0, 10) : '');
    setEditMode(true);
  };

  const handleEditSave = async (e) => {
    e.preventDefault();
    setEditSaving(true);
    setMsg('');
    try {
      const res = await fetch(`${BASE}/owner/orders/edit`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({
          id: order.id,
          address: editAddress,
          description: editDescription,
          estimated_price: editEstimatedPrice ? parseFloat(editEstimatedPrice) : null,
          planned_date: editPlannedDate || null,
        }),
      });
      if (!res.ok) throw new Error('Ошибка сохранения');
      setEditMode(false);
      setMsg('Заказ обновлён');
      await loadOrder();
    } catch (err) {
      setMsg('Ошибка: ' + err.message);
    } finally {
      setEditSaving(false);
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
        {order.request_phone && (
          <div className={styles.infoItem}>
            <span className={styles.label}>Заявка от</span>
            <span>{order.request_phone}</span>
          </div>
        )}
      </div>

      {!isCompleted && (
        <button className={styles.editBtn} onClick={openEditMode}>
          Редактировать
        </button>
      )}

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
                  <button className={styles.removeBtn} onClick={() => handleRemoveWorker(r.worker_id)} disabled={removing}>
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

      {msg && (
        <p className={msg.includes('Ошибка') ? styles.error : styles.success}>{msg}</p>
      )}

      {!isCompleted && user && reports.find((r) => r.worker_id === user.id) && (
        <>
          <h2 className={styles.sectionTitle}>Мой отчёт</h2>
          <form className={styles.reportForm} onSubmit={handleSaveReport}>
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
          </form>
        </>
      )}

      {!isCompleted && (
        <button className={styles.completeBtn} onClick={handleComplete} disabled={completing}>
          {completing ? 'Завершение...' : 'Завершить заказ'}
        </button>
      )}

      {editMode && (
        <div className={styles.overlay} onClick={() => setEditMode(false)}>
          <form className={styles.form} onClick={(e) => e.stopPropagation()} onSubmit={handleEditSave}>
            <h2>Редактировать заказ</h2>
            <label>
              Адрес
              <input type="text" value={editAddress} onChange={(e) => setEditAddress(e.target.value)} required />
            </label>
            <label>
              Описание
              <input type="text" value={editDescription} onChange={(e) => setEditDescription(e.target.value)} />
            </label>
            <label>
              Стоимость
              <input type="number" value={editEstimatedPrice} onChange={(e) => setEditEstimatedPrice(e.target.value)} />
            </label>
            <label>
              Дата
              <input type="date" value={editPlannedDate} onChange={(e) => setEditPlannedDate(e.target.value)} />
            </label>
            <div className={styles.formActions}>
              <button type="submit" className={styles.saveBtn} disabled={editSaving}>
                {editSaving ? 'Сохранение...' : 'Сохранить'}
              </button>
              <button type="button" className={styles.cancelBtn} onClick={() => setEditMode(false)}>
                Отмена
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}