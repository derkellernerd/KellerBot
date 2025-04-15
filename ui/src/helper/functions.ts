import { date, Notify } from 'quasar';

export function autoImplement<T extends object>(defaults: Partial<T> = {}) {
  return class {
    constructor(data: Partial<T> = {}) {
      Object.assign(this, defaults);
      Object.assign(this, data);
    }
  } as new (data?: T) => T;
}

export function isDateSet(value: Date) : boolean {
  if (!value) return false;
  return new Date(value) > new Date('1970-01-01')
}

export function timeFormatted(value: Date) : string {
  return date.formatDate(value, 'YYYY-MM-DD HH:mm:ss')
}

export function showSuccessToast(message: string, title?: string) {
  Notify.create({
    caption: title ?? '',
    message: message,
    type: 'info',
    timeout: 5000,
    closeBtn: true,
  })
}

export function showErrorToast(error: string, title?: string) {
  Notify.create({
    caption: title ?? '',
    message: error,
    type: 'negative',
    closeBtn: true,
  })
}
