import { afterEach, beforeEach, describe, expect, mock, test } from 'bun:test';
import { mkdtempSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import type { ResolvedAuth } from '../src/client.js';
import { listItems, type ModuleDef } from '../src/module.js';

let dir: string;
const realFetch = globalThis.fetch;

const auth: ResolvedAuth = { server: 'http://kubepi.test', token: 't' };

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });
}

function lastCall(fetchMock: ReturnType<typeof mock>): [string, RequestInit] {
  return fetchMock.mock.calls[0] as unknown as [string, RequestInit];
}

beforeEach(() => {
  dir = mkdtempSync(join(tmpdir(), 'kubepi-cli-test-'));
  process.env.KUBEPI_CONFIG_DIR = dir;
});

afterEach(() => {
  globalThis.fetch = realFetch;
  delete process.env.KUBEPI_CONFIG_DIR;
  rmSync(dir, { recursive: true, force: true });
});

describe('listItems', () => {
  const def: ModuleDef = { name: 'imagerepo', description: 'test', basePath: '/imagerepos', searchPath: '/imagerepos/search' };

  test('forceSearch module posts to searchPath even without filter/search', async () => {
    const fetchMock = mock(() => Promise.resolve(jsonResponse({ success: true, data: { items: [{ name: 'harbor' }], total: 1 } })));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const items = await listItems({ ...def, forceSearch: true }, {}, auth);
    expect(items).toEqual([{ name: 'harbor' }]);
    const [url, init] = lastCall(fetchMock);
    expect(url).toBe('http://kubepi.test/kubepi/api/v1/imagerepos/search');
    expect(init.method).toBe('POST');
    expect(JSON.parse(String(init.body))).toEqual({ conditions: {} });
  });

  test('module without forceSearch falls back to GET basePath when no conditions', async () => {
    const fetchMock = mock(() => Promise.resolve(jsonResponse({ success: true, data: [{ name: 'a' }] })));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const items = await listItems(def, {}, auth);
    expect(items).toEqual([{ name: 'a' }]);
    const [url, init] = lastCall(fetchMock);
    expect(url).toBe('http://kubepi.test/kubepi/api/v1/imagerepos');
    expect(init.method).toBe('GET');
  });

  test('search keyword pushes quick condition to searchPath', async () => {
    const fetchMock = mock(() => Promise.resolve(jsonResponse({ success: true, data: { items: [{ name: 'harbor' }], total: 1 } })));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    const items = await listItems(def, { search: 'har' }, auth);
    expect(items).toEqual([{ name: 'harbor' }]);
    const [url, init] = lastCall(fetchMock);
    expect(url).toBe('http://kubepi.test/kubepi/api/v1/imagerepos/search');
    expect(JSON.parse(String(init.body))).toEqual({ conditions: { quick: { field: 'quick', operator: 'like', value: 'har' } } });
  });
});
