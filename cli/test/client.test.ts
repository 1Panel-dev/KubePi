import { afterEach, beforeEach, describe, expect, mock, test } from 'bun:test';
import { mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { api, decodeTokenProfile, loginForToken, request, type ResolvedAuth } from '../src/client.js';
import { CliError } from '../src/errors.js';
import { loadConfig, saveConfig } from '../src/store.js';

let dir: string;
const realFetch = globalThis.fetch;

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });
}

const auth: ResolvedAuth = { server: 'http://kubepi.test', username: 'admin', password: 'kubepi' };

beforeEach(() => {
  dir = mkdtempSync(join(tmpdir(), 'kubepi-cli-test-'));
  process.env.KUBEPI_CONFIG_DIR = dir;
});

afterEach(() => {
  globalThis.fetch = realFetch;
  delete process.env.KUBEPI_CONFIG_DIR;
  rmSync(dir, { recursive: true, force: true });
});

describe('request', () => {
  test('unwraps {success,data} envelope and sends bearer token', async () => {
    saveConfig({ server: auth.server, username: 'admin', password: 'kubepi', token: 'cached-token', tokenSavedAt: Date.now() });
    const fetchMock = mock(() => Promise.resolve(jsonResponse({ success: true, data: [{ name: 'a' }] })));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const data = await api.get(auth, '/users');
    expect(data).toEqual([{ name: 'a' }]);
    const [url, init] = fetchMock.mock.calls[0] as unknown as [string, RequestInit];
    expect(url).toBe('http://kubepi.test/kubepi/api/v1/users');
    expect((init.headers as Record<string, string>).Authorization).toBe('Bearer cached-token');
  });

  test('logs in when token is stale, and retries once on 401', async () => {
    saveConfig({ server: auth.server, username: 'admin', password: 'kubepi', token: 'stale', tokenSavedAt: Date.now() - 10 * 60 * 1000 });
    const calls: string[] = [];
    globalThis.fetch = mock((input: unknown) => {
      const url = String(input);
      calls.push(url);
      if (url.endsWith('/sessions')) return Promise.resolve(new Response('fresh-token'));
      if (calls.filter((c) => c.endsWith('/users')).length === 1) {
        return Promise.resolve(jsonResponse({ success: false, code: 401, message: 'please login' }, 401));
      }
      return Promise.resolve(jsonResponse({ success: true, data: 'ok' }));
    }) as unknown as typeof fetch;

    const data = await api.get(auth, '/users');
    expect(data).toBe('ok');
    expect(calls.filter((c) => c.endsWith('/sessions')).length).toBeGreaterThanOrEqual(1);
  });

  test('maps error envelope to CliError with code', async () => {
    saveConfig({ token: 't', tokenSavedAt: Date.now() });
    globalThis.fetch = mock(() => Promise.resolve(jsonResponse({ success: false, code: 404, message: 'the server could not find the requested resource' }, 404))) as unknown as typeof fetch;
    try {
      await api.get(auth, '/users/nope');
      expect.unreachable();
    } catch (err) {
      expect(err).toBeInstanceOf(CliError);
      expect((err as CliError).code).toBe('E2002');
    }
  });

  test('raw mode returns body without unwrapping', async () => {
    saveConfig({ token: 't', tokenSavedAt: Date.now() });
    globalThis.fetch = mock(() => Promise.resolve(jsonResponse({ kind: 'PodList', items: [] }))) as unknown as typeof fetch;
    const data = (await request(auth, '/proxy/dev/k8s/api/v1/pods', { raw: true })) as { kind: string };
    expect(data.kind).toBe('PodList');
  });

  test('network failure maps to E5001', async () => {
    saveConfig({ token: 't', tokenSavedAt: Date.now() });
    globalThis.fetch = mock(() => Promise.reject(new Error('connection refused'))) as unknown as typeof fetch;
    try {
      await api.get(auth, '/users');
      expect.unreachable();
    } catch (err) {
      expect((err as CliError).code).toBe('E5001');
    }
  });

  test('loginForToken persists only the token, never the password', async () => {
    globalThis.fetch = mock(() => Promise.resolve(new Response('new-token'))) as unknown as typeof fetch;
    await loginForToken({ server: 'http://kubepi.test', username: 'admin', password: 'should-not-be-saved' });
    const cfg = loadConfig();
    expect(cfg.token).toBe('new-token');
    expect(cfg.tokenSavedAt).toBeNumber();
    expect(cfg.username).toBeUndefined();
    expect(cfg.password).toBeUndefined();
  });

  test('loginForToken surfaces MFA requirement as E1005', async () => {
    globalThis.fetch = mock(() =>
      Promise.resolve(jsonResponse({ success: false, code: 401, message: 'mfa is required' }, 401)),
    ) as unknown as typeof fetch;
    try {
      await loginForToken(auth);
      expect.unreachable();
    } catch (err) {
      expect((err as CliError).code).toBe('E1005');
    }
  });
});

describe('decodeTokenProfile', () => {
  test('decodes payload segment', () => {
    const payload = Buffer.from(JSON.stringify({ name: 'admin', isAdministrator: true })).toString('base64url');
    expect(decodeTokenProfile(`aaa.${payload}.bbb`)).toEqual({ name: 'admin', isAdministrator: true });
    expect(decodeTokenProfile('not-a-jwt')).toBeNull();
  });
});
