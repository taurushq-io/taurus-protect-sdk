/**
 * Simple key=value properties file parser.
 *
 * Supports:
 * - Lines with `key=value` format
 * - `#` comment lines (ignored)
 * - Blank lines (ignored)
 * - `\n` escape sequences in values (for PEM keys)
 */

import * as fs from 'fs';

/**
 * Parses a properties file into a key-value map.
 *
 * @param filePath - Path to the properties file
 * @returns Map of property keys to values, or null if file not found
 */
export function loadProperties(filePath: string): Map<string, string> | null {
  let content: string;
  try {
    content = fs.readFileSync(filePath, 'utf-8');
  } catch {
    return null;
  }

  const props = new Map<string, string>();
  const lines = content.split('\n');

  for (const line of lines) {
    const trimmed = line.trim();
    if (trimmed === '' || trimmed.startsWith('#')) {
      continue;
    }

    const eqIndex = trimmed.indexOf('=');
    if (eqIndex < 0) {
      continue;
    }

    const key = trimmed.substring(0, eqIndex).trim();
    let value = trimmed.substring(eqIndex + 1).trim();

    // Unescape \n in values (e.g., PEM keys)
    value = value.replace(/\\n/g, '\n');

    if (key) {
      props.set(key, value);
    }
  }

  return props;
}
