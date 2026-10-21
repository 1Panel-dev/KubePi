import * as readline from 'node:readline/promises';
import { E } from './errors.js';
import type { GlobalOpts } from './client.js';

/** Return the value when it is a non-empty string, else undefined. */
export function str(v: unknown): string | undefined {
  return typeof v === 'string' && v !== '' ? v : undefined;
}

/** Parse an optional boolean flag given as true / 'true' / 'false'. */
export function boolOpt(v: unknown): boolean | undefined {
  if (v === undefined) return undefined;
  return v === true || v === 'true';
}

/** Parse an optional numeric flag; throws E2004 when present but not a number. */
export function numOpt(v: unknown, flag: string): number | undefined {
  if (v === undefined || v === '') return undefined;
  const n = Number(v);
  if (Number.isNaN(n)) throw E.invalid(`${flag} must be a number, got: ${v}`);
  return n;
}

/** Commander collector for repeatable options. */
export function collect(value: string, prev: string[] | undefined): string[] {
  return [...(prev ?? []), value];
}

/** Pick a string option with fallback chain: opts[key] → fields[key] → fallback. */
export function pick(
  opts: Record<string, unknown>,
  fields: Record<string, unknown>,
  key: string,
  fallback?: string,
): string | undefined {
  return str(opts[key]) ?? str(fields[key]) ?? fallback;
}

/** Ask for confirmation unless --yes; refuses in non-interactive mode without --yes. */
export async function confirmAction(message: string, opts: GlobalOpts): Promise<void> {
  if (opts.yes) return;
  if (!process.stdin.isTTY) throw E.missing('--yes (this operation requires confirmation in non-interactive mode)');
  const rl = readline.createInterface({ input: process.stdin, output: process.stdout });
  try {
    const ans = await rl.question(`${message} [y/N] `);
    if (ans.trim().toLowerCase() !== 'y') {
      console.log('aborted');
      process.exit(0);
    }
  } finally {
    rl.close();
  }
}
