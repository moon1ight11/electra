import { useState } from 'react';
import styles from './Help.module.css';

const tabs = [
  {
    key: 'app',
    label: 'Приложение',
    content: (
      <>
        <h2>Как работать с приложением</h2>
        <div className={styles.block}>
          <h3>Главная</h3>
          <p>На главной вы видите ближайшие заказы. Зелёный индикатор — отчёт заполнен, оранжевый — нужно заполнить.</p>
        </div>
        <div className={styles.block}>
          <h3>В работе</h3>
          <p>Все заказы, на которые вас назначили. Нажмите на заказ, чтобы открыть подробности.</p>
        </div>
        <div className={styles.block}>
          <h3>История</h3>
          <p>Выполненные заказы. Можно посмотреть, что и когда делали.</p>
        </div>
        <div className={styles.block}>
          <h3>Отчёт</h3>
          <p>В деталях заказа заполните время, заработок, материалы и заметки. Нажмите «Сохранить отчёт».</p>
        </div>
      </>
    ),
  },
  {
    key: 'rules',
    label: 'Регламент',
    content: (
      <>
        <h2>Регламент работы</h2>
        <div className={styles.block}>
          <h3>Перед выездом</h3>
          <ul>
            <li>Проверьте адрес и описание заказа в приложении</li>
            <li>Уточните у владельца детали, если что-то непонятно</li>
            <li>Подготовьте инструмент и материалы</li>
          </ul>
        </div>
        <div className={styles.block}>
          <h3>На объекте</h3>
          <ul>
            <li>Поздоровайтесь, уточните задачу у клиента</li>
            <li>Работайте аккуратно, убирайте за собой мусор</li>
            <li>Если нужны доп. материалы — согласуйте с владельцем</li>
          </ul>
        </div>
        <div className={styles.block}>
          <h3>После работы</h3>
          <ul>
            <li>Покажите результат клиенту</li>
            <li>Заполните отчёт в приложении в тот же день</li>
            <li>Сообщите владельцу о завершении</li>
          </ul>
        </div>
      </>
    ),
  },
  {
    key: 'materials',
    label: 'Материалы',
    content: (
      <>
        <h2>Рекомендации по материалам</h2>
        <div className={styles.block}>
          <h3>Кабель</h3>
          <p>Для розеток — ВВГнг 3×2,5. Для освещения — ВВГнг 3×1,5. Для ввода — СИП 4×16.</p>
        </div>
        <div className={styles.block}>
          <h3>Автоматы</h3>
          <p>Розетки — 16А. Освещение — 10А. Вводной — 25-32А (зависит от нагрузки).</p>
        </div>
        <div className={styles.block}>
          <h3>УЗО</h3>
          <p>На розетки — 25А/30мА. На освещение — 16А/30мА.</p>
        </div>
        <div className={styles.block}>
          <h3>Розетки и выключатели</h3>
          <p>Рекомендуем Legrand, Schneider Electric, ABB. Не экономьте на качестве.</p>
        </div>
      </>
    ),
  },
  {
    key: 'safety',
    label: 'Безопасность',
    content: (
      <>
        <h2>Техника безопасности</h2>
        <div className={styles.block}>
          <h3>Общие правила</h3>
          <ul>
            <li>Всегда отключайте напряжение перед работой</li>
            <li>Проверяйте отсутствие напряжения указателем</li>
            <li>Используйте инструмент с изолированными ручками</li>
            <li>Не работайте в одиночку на высоте</li>
          </ul>
        </div>
        <div className={styles.block}>
          <h3>Средства защиты</h3>
          <ul>
            <li>Диэлектрические перчатки</li>
            <li>Защитные очки</li>
            <li>Спецобувь</li>
          </ul>
        </div>
        <div className={styles.block}>
          <h3>Экстренные номера</h3>
          <p>Пожарная — 101. Скорая — 103. Аварийная служба электросетей — 8-800-... (уточните у владельца).</p>
        </div>
      </>
    ),
  },
];

export default function Help() {
  const [activeTab, setActiveTab] = useState('app');

  const active = tabs.find((t) => t.key === activeTab);

  return (
    <div className={styles.page}>
      <h1>Помощь</h1>

      <div className={styles.tabs}>
        {tabs.map((tab) => (
          <button
            key={tab.key}
            className={activeTab === tab.key ? styles.tabActive : styles.tab}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className={styles.content}>
        {active && active.content}
      </div>
    </div>
  );
}