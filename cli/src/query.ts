import { E } from './errors.js';
import { getPath } from './output.js';
import { numOpt } from './shared.js';

export type FilterOp = ':' | '=' | '!=' | '>' | '<' | '>=' | '<=' | '~' | '!~';

export interface FilterTerm {
  field: string;
  op: FilterOp;
  value: string;
}

/** OR groups; each group is a list of AND terms. */
export type FilterGroups = FilterTerm[][];

const OP_PATTERN = /^(.*?)(!~|!=|>=|<=|>|<|~|=|:)(.*)$/;

export function parseFilters(inputs?: string[]): FilterGroups {
  if (!inputs || inputs.length === 0) return [];
  return inputs.map((input) =>
    input.split(',').map((term) => {
      const m = OP_PATTERN.exec(term.trim());
      if (!m || !m[1]) throw E.badRequest(`invalid filter expression: ${term}`);
      return { field: m[1].trim(), op: m[2] as FilterOp, value: m[3].trim() };
    }),
  );
}

function comparable(v: unknown): string | number {
  if (typeof v === 'number') return v;
  if (v === null || v === undefined) return '';
  return String(v);
}

function matchTerm(item: unknown, term: FilterTerm): boolean {
  const raw = getPath(item, term.field);
  const left = comparable(raw);
  const right = term.value;
  const leftNum = typeof left === 'number' ? left : Number(left);
  const rightNum = Number(right);
  const numeric = !Number.isNaN(leftNum) && !Number.isNaN(rightNum) && right !== '';
  switch (term.op) {
    case ':':
    case '=':
      return Array.isArray(raw) ? raw.map(String).includes(right) : String(left) === right;
    case '!=':
      return Array.isArray(raw) ? !raw.map(String).includes(right) : String(left) !== right;
    case '>':
      return numeric ? leftNum > rightNum : String(left) > right;
    case '<':
      return numeric ? leftNum < rightNum : String(left) < right;
    case '>=':
      return numeric ? leftNum >= rightNum : String(left) >= right;
    case '<=':
      return numeric ? leftNum <= rightNum : String(left) <= right;
    case '~':
      return String(left).toLowerCase().includes(right.toLowerCase());
    case '!~':
      return !String(left).toLowerCase().includes(right.toLowerCase());
  }
}

export function matchGroups(item: unknown, groups: FilterGroups): boolean {
  if (groups.length === 0) return true;
  return groups.some((andTerms) => andTerms.every((t) => matchTerm(item, t)));
}

function collectStrings(v: unknown, depth: number, out: string[]): void {
  if (depth > 3 || v === null || v === undefined) return;
  if (typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean') {
    out.push(String(v));
    return;
  }
  if (Array.isArray(v)) {
    for (const item of v.slice(0, 50)) collectStrings(item, depth + 1, out);
    return;
  }
  if (typeof v === 'object') {
    for (const val of Object.values(v)) collectStrings(val, depth + 1, out);
  }
}

/** Client-side fuzzy search across all scalar fields of an item. */
export function matchSearch(item: unknown, keyword: string): boolean {
  const strings: string[] = [];
  collectStrings(item, 0, strings);
  const kw = keyword.toLowerCase();
  return strings.some((s) => s.toLowerCase().includes(kw));
}

export function applySort(items: unknown[], sort?: string): unknown[] {
  if (!sort) return items;
  const keys = sort.split(',').map((k) => k.trim()).filter(Boolean);
  const sorted = [...items];
  sorted.sort((a, b) => {
    for (const key of keys) {
      const desc = key.endsWith('_desc');
      const field = key.replace(/_(asc|desc)$/, '');
      const va = comparable(getPath(a, field));
      const vb = comparable(getPath(b, field));
      let cmp: number;
      const na = Number(va);
      const nb = Number(vb);
      if (!Number.isNaN(na) && !Number.isNaN(nb) && va !== '' && vb !== '') cmp = na - nb;
      else cmp = String(va).localeCompare(String(vb));
      if (cmp !== 0) return desc ? -cmp : cmp;
    }
    return 0;
  });
  return sorted;
}

interface Condition {
  field: string;
  operator: string;
  value: string;
}

export interface ClientQueryOpts {
  filter?: unknown;
  search?: string;
  sort?: string;
  limit?: unknown;
}

/**
 * Client-side post-processing for list-like results: --filter, --search,
 * --sort and --limit, in this order. Shared by the plain list path and
 * custom actions returning arrays; conditions already pushed server-side
 * are re-applied idempotently.
 */
export function applyClientQuery(items: unknown[], opts: ClientQueryOpts): unknown[] {
  const groups = parseFilters(Array.isArray(opts.filter) ? (opts.filter as string[]) : undefined);
  let out = items.filter((i) => matchGroups(i, groups));
  if (opts.search) out = out.filter((i) => matchSearch(i, opts.search!));
  out = applySort(out, opts.sort);
  const limit = numOpt(opts.limit, '--limit');
  if (limit !== undefined) out = out.slice(0, limit);
  return out;
}

/**
 * Build KubePi /search conditions from CLI filters.
 * Server supports eq / ne / like / not like plus a special `quick` fuzzy key.
 * OR semantics cannot be expressed server-side, so multiple --filter groups
 * disable server filtering (client-side still applies).
 */
export function buildConditions(groups: FilterGroups, search?: string): Record<string, Condition> | null {
  const conditions: Record<string, Condition> = {};
  if (search) conditions.quick = { field: 'quick', operator: 'like', value: search };
  if (groups.length > 1) return Object.keys(conditions).length > 0 ? conditions : search ? conditions : null;
  for (const term of groups[0] ?? []) {
    if (term.field.includes('.')) return search ? conditions : null;
    switch (term.op) {
      case ':':
      case '=':
        conditions[term.field] = { field: term.field, operator: 'eq', value: term.value };
        break;
      case '!=':
        conditions[term.field] = { field: term.field, operator: 'ne', value: term.value };
        break;
      case '~':
        conditions[term.field] = { field: term.field, operator: 'like', value: term.value };
        break;
      case '!~':
        conditions[term.field] = { field: term.field, operator: 'not like', value: term.value };
        break;
      default:
        return Object.keys(conditions).length > 0 ? conditions : null;
    }
  }
  return Object.keys(conditions).length > 0 ? conditions : null;
}
