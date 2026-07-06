import { useState, useRef } from 'react';
import { Link } from 'react-router-dom';
import styles from './Landing.module.css';
import logo from '../../logo.svg';

const works = [
  {
    id: 1,
    title: 'Подключение дома',
    description: 'Полная электрификация частного дома 150 м²: ввод, щиток, розетки, освещение',
    price: 'от 25 000 ₽',
    image: '/images/works/work-1.jpg',
  },
  {
    id: 2,
    title: 'Замена проводки',
    description: 'Замена старой алюминиевой проводки на медную в квартире-студии',
    price: 'от 18 000 ₽',
    image: '/images/works/work-2.jpg',
  },
  {
    id: 3,
    title: 'Электрика в бане',
    description: 'Проводка, освещение, тёплый пол, подключение печи',
    price: 'от 15 000 ₽',
    image: '/images/works/work-3.jpg',
  },
  {
    id: 4,
    title: 'Щиток под ключ',
    description: 'Сборка и установка распределительного щита с автоматами и УЗО',
    price: 'от 8 000 ₽',
    image: '/images/works/work-4.jpg',
  },
  {
    id: 5,
    title: 'Освещение участка',
    description: 'Установка уличных светильников, прокладка кабеля, автоматика',
    price: 'от 12 000 ₽',
    image: '/images/works/work-5.jpg',
  },
  {
    id: 6,
    title: 'Диагностика сети',
    description: 'Поиск неисправностей, замеры, проверка заземления и автоматов',
    price: 'от 3 000 ₽',
    image: '/images/works/work-6.jpg',
  },
];

export default function Landing() {
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [comment, setComment] = useState('');
  const [sent, setSent] = useState(false);
  const [error, setError] = useState('');
  const [scrollPos, setScrollPos] = useState(0);
  const trackRef = useRef(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch('/api/v1/public/requests', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, phone, comment }),
      });
      if (!res.ok) throw new Error('Ошибка отправки');
      setSent(true);
    } catch (err) {
      setError('Не удалось отправить заявку. Попробуйте позже.');
    }
  };

  const scroll = (direction) => {
    const track = trackRef.current;
    if (!track) return;
    const cardWidth = track.querySelector(`.${styles.workCard}`)?.offsetWidth || 300;
    const gap = 24;
    const step = cardWidth + gap;
    const newPos = direction === 'left' ? scrollPos - step : scrollPos + step;
    const maxScroll = track.scrollWidth - track.offsetWidth;
    const clamped = Math.max(0, Math.min(newPos, maxScroll));
    track.style.transform = `translateX(-${clamped}px)`;
    setScrollPos(clamped);
  };

  return (
    <div className={styles.page}>
      {/* Шапка */}
      <header className={styles.header}>
        <img src={logo} alt="Electra" className={styles.logoImg} />
        <Link to="/login" className={styles.loginBtn}>Вход</Link>
      </header>

      {/* Главный блок */}
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <h1>Электрика для вашего дома</h1>
          <p>Подключение, замена проводки, щитки — всё, от розетки до ввода. Работаем в частных домах и квартирах.</p>
          <div className={styles.contacts}>
            <span>📞 +7 (999) 123-45-67</span>
            <span>✉️ electra@example.com</span>
          </div>
        </div>
      </section>

      {/* Примеры работ — карусель */}
      <section className={styles.works}>
        <h2>Примеры работ</h2>
        <div className={styles.carousel}>
          <button className={styles.arrow} onClick={() => scroll('left')} aria-label="Назад">‹</button>
          <div className={styles.trackWrapper}>
            <div className={styles.track} ref={trackRef}>
              {works.map((w) => (
                <div key={w.id} className={styles.workCard}>
                  <div className={styles.workImage}>
                    <img src={w.image} alt={w.title} />
                  </div>
                  <h3>{w.title}</h3>
                  <p>{w.description}</p>
                  <span className={styles.price}>{w.price}</span>
                </div>
              ))}
            </div>
          </div>
          <button className={styles.arrow} onClick={() => scroll('right')} aria-label="Вперёд">›</button>
        </div>
      </section>

      {/* Почему мы */}
      <section className={styles.why}>
        <h2>Почему мы</h2>
        <div className={styles.whyGrid}>
          <div className={styles.whyCard}>
            <div className={styles.whyIcon}>
              <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
                <rect x="10" y="8" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="19" y="8" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="28" y="8" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="10" y="15" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="19" y="15" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="28" y="15" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="10" y="22" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="19" y="22" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
                <rect x="28" y="22" width="8" height="6" rx="1" stroke="#3D4A3E" strokeWidth="2" />
              </svg>
            </div>
            <h3>Работаем с 2020</h3>
            <p>Более 200 выполненных объектов — от квартир до коттеджей</p>
          </div>
          <div className={styles.whyCard}>
            <div className={styles.whyIcon}>
              <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
                <path d="M20 6L23.5 14L32 15L26 20.5L27.5 29L20 24.5L12.5 29L14 20.5L8 15L16.5 14L20 6Z" fill="#3D4A3E" />
              </svg>
            </div>
            <h3>Даём гарантию</h3>
            <p>12 месяцев на все виды работ и материалы</p>
          </div>
          <div className={styles.whyCard}>
            <div className={styles.whyIcon}>
              <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
                <circle cx="20" cy="20" r="14" stroke="#3D4A3E" strokeWidth="2" />
                <line x1="20" y1="12" x2="20" y2="20" stroke="#3D4A3E" strokeWidth="2.5" strokeLinecap="round" />
                <line x1="20" y1="20" x2="26" y2="16" stroke="#D4A030" strokeWidth="2" strokeLinecap="round" />
                <circle cx="20" cy="20" r="2" fill="#3D4A3E" />
              </svg>
            </div>
            <h3>Выезжаем быстро</h3>
            <p>Осмотр и оценка в день звонка или на следующий день</p>
          </div>
          <div className={styles.whyCard}>
            <div className={styles.whyIcon}>
              <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
                <rect x="10" y="8" width="20" height="24" rx="2" stroke="#3D4A3E" strokeWidth="2.5" />
                <line x1="10" y1="15" x2="30" y2="15" stroke="#3D4A3E" strokeWidth="2" />
                <line x1="14" y1="22" x2="26" y2="22" stroke="#D4A030" strokeWidth="2" strokeLinecap="round" />
              </svg>
            </div>
            <h3>Работаем по договору</h3>
            <p>Фиксируем стоимость до начала работ, никаких сюрпризов</p>
          </div>
        </div>
      </section>

      {/* Форма заявки */}
      <section className={styles.request}>
        <h2>Оставьте заявку</h2>
        <p>Опишите задачу, и мы перезвоним в течение часа</p>
        {sent ? (
          <div className={styles.success}>Спасибо! Мы скоро свяжемся с вами.</div>
        ) : (
          <form className={styles.form} onSubmit={handleSubmit}>
            {error && <div className={styles.formError}>{error}</div>}
            <input type="text" placeholder="Ваше имя" value={name} onChange={(e) => setName(e.target.value)} required />
            <input type="text" placeholder="Телефон" value={phone} onChange={(e) => setPhone(e.target.value)} required />
            <textarea placeholder="Что нужно сделать?" value={comment} onChange={(e) => setComment(e.target.value)} rows={3} />
            <button type="submit">Отправить</button>
          </form>
        )}
      </section>

      {/* FAQ */}
      <section className={styles.faq}>
        <h2>Частые вопросы</h2>
        <div className={styles.faqList}>
          <details className={styles.faqItem}>
            <summary>Сколько стоит выезд на оценку?</summary>
            <p>Выезд в пределах города — бесплатно. В область — от 500 ₽, которые вычитаются из стоимости работ при заказе.</p>
          </details>
          <details className={styles.faqItem}>
            <summary>Делаете ли вы «под ключ»?</summary>
            <p>Да, от проекта до сдачи. Сами закупаем материалы, делаем монтаж, убираем за собой.</p>
          </details>
          <details className={styles.faqItem}>
            <summary>Можно ли вызвать электрика в выходной?</summary>
            <p>Да, работаем без выходных. Срочный вызов — в течение 2 часов.</p>
          </details>
          <details className={styles.faqItem}>
            <summary>Какие материалы используете?</summary>
            <p>Работаем с проверенными брендами: ABB, Legrand, Schneider Electric. Чеки и паспорта на всё предоставляем.</p>
          </details>
        </div>
      </section>

      {/* Подвал */}
      <footer className={styles.footer}>
        <span>Electra</span>
        <span>Работаем с 2020 года</span>
        <span>📞 +7 (999) 123-45-67</span>
      </footer>
    </div>
  );
}