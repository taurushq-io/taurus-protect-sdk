/**
 * Safe conversion utilities for mapping OpenAPI DTOs to domain models.
 * These functions handle null/undefined values gracefully.
 */

/**
 * Safely converts a value to a string.
 */
export function safeString(value: unknown): string | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  return String(value);
}

/**
 * Safely converts a value to a string, with a default value.
 */
export function safeStringDefault(value: unknown, defaultValue: string): string {
  if (value === null || value === undefined) {
    return defaultValue;
  }
  return String(value);
}

/**
 * Safely converts a value to a number.
 */
export function safeInt(value: unknown): number | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  const num = typeof value === 'number' ? value : parseInt(String(value), 10);
  return isNaN(num) ? undefined : num;
}

/**
 * Safely converts a value to a number, with a default value.
 */
export function safeIntDefault(value: unknown, defaultValue: number): number {
  const result = safeInt(value);
  return result === undefined ? defaultValue : result;
}

/**
 * Safely converts a value to a float.
 */
export function safeFloat(value: unknown): number | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  const num = typeof value === 'number' ? value : parseFloat(String(value));
  return isNaN(num) ? undefined : num;
}

/**
 * Safely converts a value to a boolean.
 */
export function safeBool(value: unknown): boolean | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  if (typeof value === 'boolean') {
    return value;
  }
  if (typeof value === 'string') {
    const lower = value.toLowerCase();
    if (lower === 'true' || lower === '1' || lower === 'yes') return true;
    if (lower === 'false' || lower === '0' || lower === 'no') return false;
  }
  if (typeof value === 'number') {
    return value !== 0;
  }
  return undefined;
}

/**
 * Safely converts a value to a boolean, with a default value.
 */
export function safeBoolDefault(value: unknown, defaultValue: boolean): boolean {
  const result = safeBool(value);
  return result === undefined ? defaultValue : result;
}

/**
 * Safely converts a value to a Date.
 */
export function safeDate(value: unknown): Date | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  if (value instanceof Date) {
    return value;
  }
  const date = new Date(String(value));
  return isNaN(date.getTime()) ? undefined : date;
}

/**
 * Safely converts a value to an ISO date string.
 */
export function safeDateString(value: unknown): string | undefined {
  const date = safeDate(value);
  return date?.toISOString();
}

/**
 * Maps an array safely, filtering out undefined results.
 */
export function safeMap<T, U>(
  array: T[] | null | undefined,
  mapper: (item: T) => U | undefined
): U[] {
  if (!array) {
    return [];
  }
  const result: U[] = [];
  for (const item of array) {
    const mapped = mapper(item);
    if (mapped !== undefined) {
      result.push(mapped);
    }
  }
  return result;
}
