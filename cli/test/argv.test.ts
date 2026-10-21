import { describe, expect, test } from 'bun:test';
import { hoistGlobalOptions } from '../src/argv.js';

describe('hoistGlobalOptions', () => {
  test('moves global flags before module args', () => {
    const out = hoistGlobalOptions(['node', 'kubepi', 'user', '--format=json', 'admin']);
    expect(out).toEqual(['node', 'kubepi', '--format=json', 'user', 'admin']);
  });

  test('handles --flag value form', () => {
    const out = hoistGlobalOptions(['node', 'kubepi', 'user', 'create', '--name', 'x', '--filter', 'a:1']);
    expect(out).toEqual(['node', 'kubepi', '--filter', 'a:1', 'user', 'create', '--name', 'x']);
  });

  test('bool flags and unknown flags', () => {
    const out = hoistGlobalOptions(['node', 'kubepi', 'user', 'delete', 'tom', '--yes', '--nick-name=t']);
    expect(out).toEqual(['node', 'kubepi', '--yes', 'user', 'delete', 'tom', '--nick-name=t']);
  });

  test('keeps order stable for non-global flags', () => {
    const out = hoistGlobalOptions(['node', 'kubepi', 'cluster', 'create', '--name=dev', '--mode=bearer']);
    expect(out).toEqual(['node', 'kubepi', 'cluster', 'create', '--name=dev', '--mode=bearer']);
  });

  test('does not eat the next flag as a value', () => {
    const out = hoistGlobalOptions(['node', 'kubepi', 'user', '--search', '--format=json']);
    expect(out).toEqual(['node', 'kubepi', '--search', '--format=json', 'user']);
  });

  test('missing value at end stays bare', () => {
    const out = hoistGlobalOptions(['node', 'kubepi', 'user', '--search']);
    expect(out).toEqual(['node', 'kubepi', '--search', 'user']);
  });
});
