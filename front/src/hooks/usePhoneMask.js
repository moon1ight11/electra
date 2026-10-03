import { useState, useCallback } from 'react';

export function usePhoneMask(initial = '') {
  const [value, setValue] = useState(initial);

  const onChange = (e) => {
    const raw = e.target.value.replace(/\D/g, '');
    let formatted = '+7 (';

    if (raw.length > 1) formatted += raw.slice(1, 4);
    if (raw.length > 4) formatted += ') ' + raw.slice(4, 7);
    if (raw.length > 7) formatted += '-' + raw.slice(7, 9);
    if (raw.length > 9) formatted += '-' + raw.slice(9, 11);

    setValue(formatted);
  };

  const reset = useCallback(() => setValue(''), []);

  return { value, onChange, reset };
}