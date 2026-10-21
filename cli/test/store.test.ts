import { afterEach, beforeEach, describe, expect, test } from 'bun:test';
import { mkdtempSync, rmSync, statSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { clearConfig, configFile, loadConfig, patchConfig, saveConfig } from '../src/store.js';

let dir: string;

beforeEach(() => {
  dir = mkdtempSync(join(tmpdir(), 'kubepi-cli-test-'));
  process.env.KUBEPI_CONFIG_DIR = dir;
});

afterEach(() => {
  delete process.env.KUBEPI_CONFIG_DIR;
  rmSync(dir, { recursive: true, force: true });
});

describe('store', () => {
  test('loadConfig returns {} when missing', () => {
    expect(loadConfig()).toEqual({});
  });

  test('save/load roundtrip with 0600 permission', () => {
    saveConfig({ server: 'http://x', username: 'u', password: 'p', token: 't', tokenSavedAt: 1 });
    expect(loadConfig()).toEqual({ server: 'http://x', username: 'u', password: 'p', token: 't', tokenSavedAt: 1 });
    const mode = statSync(configFile()).mode & 0o777;
    expect(mode).toBe(0o600);
  });

  test('patchConfig merges', () => {
    saveConfig({ server: 'http://x', username: 'u' });
    patchConfig({ token: 't2' });
    expect(loadConfig()).toEqual({ server: 'http://x', username: 'u', token: 't2' });
  });

  test('clearConfig removes the file', () => {
    saveConfig({ server: 'http://x' });
    clearConfig();
    expect(loadConfig()).toEqual({});
  });

  test('corrupt file falls back to {}', () => {
    saveConfig({ server: 'http://x' });
    writeFileSync(configFile(), 'not-json');
    expect(loadConfig()).toEqual({});
  });
});
