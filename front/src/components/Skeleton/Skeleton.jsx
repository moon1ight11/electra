import styles from './Skeleton.module.css';

export default function Skeleton({ count = 3, height = 72 }) {
  return (
    <div className={styles.wrapper}>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className={styles.card} style={{ height }}>
          <div className={styles.line} style={{ width: '60%' }} />
          <div className={styles.line} style={{ width: '40%' }} />
        </div>
      ))}
    </div>
  );
}