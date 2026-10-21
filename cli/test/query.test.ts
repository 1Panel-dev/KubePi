import { describe, expect, test } from 'bun:test';
import { applyClientQuery, applySort, buildConditions, matchGroups, matchSearch, parseFilters } from '../src/query.js';

describe('parseFilters', () => {
  test('parses operators', () => {
    const groups = parseFilters(['name:admin,age>18']);
    expect(groups).toHaveLength(1);
    expect(groups[0]).toEqual([
      { field: 'name', op: ':', value: 'admin' },
      { field: 'age', op: '>', value: '18' },
    ]);
  });

  test('multiple --filter means OR', () => {
    const groups = parseFilters(['status:active', 'status:pending']);
    expect(groups).toHaveLength(2);
  });

  test('operator precedence: != and !~ before > < : ~ =', () => {
    expect(parseFilters(['a!=1'])[0][0].op).toBe('!=');
    expect(parseFilters(['a!~x'])[0][0].op).toBe('!~');
    expect(parseFilters(['a>=1'])[0][0].op).toBe('>=');
    expect(parseFilters(['a<=1'])[0][0].op).toBe('<=');
  });

  test('invalid filter throws', () => {
    expect(() => parseFilters(['>x'])).toThrow();
  });
});

describe('matchGroups', () => {
  const item = { name: 'admin', age: 30, status: 'active', labels: ['a', 'b'] };

  test('empty groups match everything', () => {
    expect(matchGroups(item, [])).toBe(true);
  });

  test('AND within group, OR across groups', () => {
    expect(matchGroups(item, parseFilters(['name:admin,age>20']))).toBe(true);
    expect(matchGroups(item, parseFilters(['name:admin,age>40']))).toBe(false);
    expect(matchGroups(item, parseFilters(['name:tom', 'age:30']))).toBe(true);
  });

  test('operators', () => {
    expect(matchGroups(item, parseFilters(['age>=30']))).toBe(true);
    expect(matchGroups(item, parseFilters(['age<30']))).toBe(false);
    expect(matchGroups(item, parseFilters(['name!=admin']))).toBe(false);
    expect(matchGroups(item, parseFilters(['name~dmi']))).toBe(true);
    expect(matchGroups(item, parseFilters(['name!~dmi']))).toBe(false);
    expect(matchGroups(item, parseFilters(['labels:a']))).toBe(true);
    expect(matchGroups(item, parseFilters(['labels:c']))).toBe(false);
  });

  test('dot path fields', () => {
    const nested = { metadata: { name: 'pod-a' } };
    expect(matchGroups(nested, parseFilters(['metadata.name:pod-a']))).toBe(true);
  });
});

describe('matchSearch', () => {
  test('deep fuzzy match', () => {
    expect(matchSearch({ a: { b: ['hello world'] } }, 'WORLD')).toBe(true);
    expect(matchSearch({ a: 1 }, 'nope')).toBe(false);
  });
});

describe('applySort', () => {
  test('sorts asc/desc with numeric awareness', () => {
    const items: Record<string, unknown>[] = [{ name: 'b', n: 2 }, { name: 'a', n: 10 }, { name: 'c', n: 1 }];
    expect((applySort(items, 'n_asc') as typeof items).map((i) => i.n)).toEqual([1, 2, 10]);
    expect((applySort(items, 'n_desc') as typeof items).map((i) => i.n)).toEqual([10, 2, 1]);
    expect((applySort(items, 'name_asc') as typeof items).map((i) => i.name)).toEqual(['a', 'b', 'c']);
  });
});

describe('applyClientQuery', () => {
  const items = [
    { name: 'alpha', n: 3, status: 'active' },
    { name: 'beta', n: 1, status: 'inactive' },
    { name: 'gamma', n: 2, status: 'active' },
  ];

  test('no options returns items unchanged', () => {
    expect(applyClientQuery(items, {})).toEqual(items);
  });

  test('applies limit', () => {
    expect(applyClientQuery(items, { limit: '2' })).toHaveLength(2);
  });

  test('applies filter, search, sort and limit in order', () => {
    const out = applyClientQuery(items, {
      filter: ['status:active'],
      search: 'a',
      sort: 'n_asc',
      limit: '1',
    }) as typeof items;
    expect(out.map((i) => i.name)).toEqual(['gamma']);
  });

  test('non-array filter is ignored (action-local --filter overrides)', () => {
    expect(applyClientQuery(items, { filter: '(uid=%s)' })).toEqual(items);
  });

  test('invalid limit throws E2004', () => {
    expect(() => applyClientQuery(items, { limit: 'abc' })).toThrow(/must be a number/);
  });
});

describe('buildConditions', () => {
  test('maps operators to server conditions', () => {
    const conditions = buildConditions(parseFilters(['name:admin,type~local,age!=3']));
    expect(conditions).toEqual({
      name: { field: 'name', operator: 'eq', value: 'admin' },
      type: { field: 'type', operator: 'like', value: 'local' },
      age: { field: 'age', operator: 'ne', value: '3' },
    });
  });

  test('search becomes quick condition', () => {
    expect(buildConditions([], 'foo')).toEqual({ quick: { field: 'quick', operator: 'like', value: 'foo' } });
  });

  test('OR groups or dot paths fall back to quick-only / null', () => {
    expect(buildConditions(parseFilters(['a:1', 'b:2']))).toBeNull();
    expect(buildConditions(parseFilters(['a.b:1']))).toBeNull();
    expect(buildConditions(parseFilters(['a:1', 'b:2']), 'x')).toEqual({
      quick: { field: 'quick', operator: 'like', value: 'x' },
    });
  });
});
