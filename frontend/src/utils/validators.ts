export function isEmail(value: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
}

export function required(value: string | undefined | null, label: string): string | null {
  if (!value?.trim()) return `${label} is required`;
  return null;
}
