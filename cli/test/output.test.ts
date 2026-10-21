import { describe, expect, test } from 'bun:test';
import { getPath, markdownTable, pickFields, setPath } from '../src/output.js';

describe('getPath/setPath', () => {
  test('dot paths', () => {
    const obj: Record<string, unknown> = {};
    setPath(obj, 'a.b.c', 1);
    expect(obj).toEqual({ a: { b: { c: 1 } } });
    expect(getPath(obj, 'a.b.c')).toBe(1);
    expect(getPath(obj, 'a.x.y')).toBeUndefined();
    expect(getPath(null, 'a')).toBeUndefined();
  });
});

describe('pickFields', () => {
  test('picks nested fields', () => {
    const item = { name: 'a', spec: { connect: { forward: { apiServer: 'x' } } }, other: 1 };
    expect(pickFields(item, ['name', 'spec.connect.forward.apiServer'])).toEqual({
      name: 'a',
      spec: { connect: { forward: { apiServer: 'x' } } },
    });
  });
});

describe('markdownTable', () => {
  test('renders header, separator and rows', () => {
    const table = markdownTable([{ a: 1, b: 'x|y' }], ['a', 'b']);
    const lines = table.split('\n');
    expect(lines[0]).toBe('| a | b |');
    expect(lines[1]).toBe('| --- | --- |');
    expect(lines[2]).toBe('| 1 | x\\|y |');
  });

  test('empty list', () => {
    expect(markdownTable([], ['a'])).toBe('(empty)');
  });

  test('objects are stringified, newlines flattened', () => {
    const table = markdownTable([{ a: { x: 1 }, b: 'l1\nl2' }], ['a', 'b']);
    expect(table.split('\n')[2]).toBe('| {"x":1} | l1 l2 |');
  });
});
