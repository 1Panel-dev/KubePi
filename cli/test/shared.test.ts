import { describe, expect, test } from 'bun:test';
import { CliError } from '../src/errors.js';
import { boolOpt, collect, numOpt, pick, str } from '../src/shared.js';

describe('str', () => {
  test('non-empty string passes through', () => {
    expect(str('x')).toBe('x');
    expect(str('')).toBeUndefined();
    expect(str(1)).toBeUndefined();
    expect(str(undefined)).toBeUndefined();
  });
});

describe('boolOpt', () => {
  test('parses booleans', () => {
    expect(boolOpt(undefined)).toBeUndefined();
    expect(boolOpt(true)).toBe(true);
    expect(boolOpt('true')).toBe(true);
    expect(boolOpt('false')).toBe(false);
    expect(boolOpt('x')).toBe(false);
  });
});

describe('numOpt', () => {
  test('parses numbers', () => {
    expect(numOpt(undefined, '--limit')).toBeUndefined();
    expect(numOpt('', '--limit')).toBeUndefined();
    expect(numOpt('10', '--limit')).toBe(10);
  });

  test('NaN throws E2004 with flag name', () => {
    try {
      numOpt('abc', '--limit');
      expect.unreachable();
    } catch (err) {
      expect(err).toBeInstanceOf(CliError);
      expect((err as CliError).code).toBe('E2004');
      expect((err as CliError).message).toContain('--limit');
    }
  });
});

describe('collect', () => {
  test('accumulates repeatable options', () => {
    expect(collect('a', undefined)).toEqual(['a']);
    expect(collect('b', ['a'])).toEqual(['a', 'b']);
  });
});

describe('pick', () => {
  test('opts wins over fields, fallback last', () => {
    expect(pick({ name: 'o' }, { name: 'f' }, 'name')).toBe('o');
    expect(pick({}, { name: 'f' }, 'name')).toBe('f');
    expect(pick({}, {}, 'name', 'd')).toBe('d');
    expect(pick({}, {}, 'name')).toBeUndefined();
  });
});
