import type { GlobalOpts } from './client.js';

export function getPath(obj: unknown, path: string): unknown {
  let cur = obj;
  for (const key of path.split('.')) {
    if (cur === null || cur === undefined || typeof cur !== 'object') return undefined;
    cur = (cur as Record<string, unknown>)[key];
  }
  return cur;
}

export function setPath(obj: Record<string, unknown>, path: string, value: unknown): void {
  const keys = path.split('.');
  let cur = obj;
  for (let i = 0; i < keys.length - 1; i++) {
    const k = keys[i];
    if (typeof cur[k] !== 'object' || cur[k] === null || Array.isArray(cur[k])) cur[k] = {};
    cur = cur[k] as Record<string, unknown>;
  }
  cur[keys[keys.length - 1]] = value;
}

export function pickFields(item: unknown, paths: string[]): Record<string, unknown> {
  const out: Record<string, unknown> = {};
  for (const p of paths) setPath(out, p, getPath(item, p));
  return out;
}

function cellText(v: unknown): string {
  if (v === null || v === undefined) return '';
  let s: string;
  if (typeof v === 'object') s = JSON.stringify(v);
  else s = String(v);
  return s.replace(/\r?\n/g, ' ').replace(/\|/g, '\\|');
}

function inferColumns(items: unknown[], cap = 8): string[] {
  const cols: string[] = [];
  for (const item of items.slice(0, 20)) {
    if (!item || typeof item !== 'object' || Array.isArray(item)) continue;
    for (const [k, v] of Object.entries(item as Record<string, unknown>)) {
      if (cols.includes(k)) continue;
      if (typeof v === 'object' && v !== null) continue;
      cols.push(k);
      if (cols.length >= cap) return cols;
    }
  }
  return cols;
}

export function markdownTable(items: unknown[], columns: string[]): string {
  if (items.length === 0) return '(empty)';
  const header = `| ${columns.join(' | ')} |`;
  const sep = `| ${columns.map(() => '---').join(' | ')} |`;
  const rows = items.map((item) => `| ${columns.map((c) => cellText(getPath(item, c))).join(' | ')} |`);
  return [header, sep, ...rows].join('\n');
}

function kvList(obj: Record<string, unknown>): string {
  return Object.entries(obj)
    .map(([k, v]) => `- ${k}: ${cellText(v)}`)
    .join('\n');
}

export interface PrintContext {
  opts: GlobalOpts;
  defaultColumns?: string[];
}

/** Print command result: list as markdown table, object as k-v list, raw JSON with --format=json. */
export function printData(data: unknown, ctx: PrintContext): void {
  const { opts } = ctx;
  const pick = opts.pick?.split(',').map((s) => s.trim()).filter(Boolean);

  let out = data;
  if (pick && out !== null && out !== undefined) {
    out = Array.isArray(out) ? out.map((i) => pickFields(i, pick)) : pickFields(out, pick);
  }

  if (opts.format === 'json') {
    console.log(JSON.stringify(out, null, 2));
    return;
  }
  if (out === null || out === undefined) {
    console.log('ok');
    return;
  }
  if (Array.isArray(out)) {
    const columns = pick ?? ctx.defaultColumns ?? inferColumns(out);
    console.log(markdownTable(out, columns));
    return;
  }
  if (typeof out === 'object') {
    console.log(kvList(out as Record<string, unknown>));
    return;
  }
  console.log(String(out));
}
