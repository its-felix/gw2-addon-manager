import { useMemo } from 'react';
import { DateFormat } from '../../../lib/preferences.model';
import { usePreferences } from './use-preferences';

interface DateFormatter {
  formatDate: (v: string | Date) => string;
  formatTime: (v: string | Date) => string;
  formatDateTime: (v: string | Date) => string;
}

export const ISO8601DateFormatter: DateFormatter = {
  formatDate: (v) => safeDate(v).toISOString().split('T')[0],
  formatTime: (v) => safeDate(v).toISOString().split('T')[1],
  formatDateTime: (v) => safeDate(v).toISOString(),
} as const;

export const SystemDateFormatter: DateFormatter = {
  formatDate: (v) => safeDate(v).toLocaleDateString(),
  formatTime: (v) => safeDate(v).toLocaleTimeString(),
  formatDateTime: (v) => safeDate(v).toLocaleString(),
} as const;

export function useDateFormat() {
  const [preferences] = usePreferences();

  return useMemo<DateFormatter>(() => {
    switch (preferences.dateFormat) {
      case DateFormat.ISO_8601:
        return ISO8601DateFormatter;

      case DateFormat.SYSTEM:
      default:
        return SystemDateFormatter;
    }
  }, [preferences]);
}

function safeDate(v: string | Date) {
  if (typeof v === 'string') {
    return new Date(v);
  }

  return v;
}
