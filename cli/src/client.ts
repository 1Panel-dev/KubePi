import { CliError, E, errorMessage } from './errors.js';
import { loadConfig, patchConfig } from './store.js';

export const API_PREFIX = '/kubepi/api/v1';

/** JWT issued by KubePi lives 10 minutes by default; refresh a bit earlier. */
const TOKEN_MAX_AGE_MS = 9 * 60 * 1000;

export interface GlobalOpts {
  server?: string;
  username?: string;
  password?: string;
  token?: string;
  format?: string;
  pick?: string;
  filter?: string[];
  search?: string;
  sort?: string;
  page?: string;
  pageSize?: string;
  all?: boolean;
  limit?: string;
  yes?: boolean;
  data?: string;
}

export interface ResolvedAuth {
  server: string;
  username?: string;
  password?: string;
  /** Explicit token from --token / KUBEPI_TOKEN; never refreshed automatically. */
  token?: string;
}

/** Resolve credentials with precedence: command-line options > env vars > config file. */
export function resolveAuth(opts: GlobalOpts): ResolvedAuth {
  assertNoSwallowedFlags(opts);
  const cfg = loadConfig();
  const server = opts.server || process.env.KUBEPI_URL || cfg.server;
  if (!server) throw E.notLoggedIn();
  return {
    server: server.replace(/\/+$/, ''),
    username: opts.username || process.env.KUBEPI_USERNAME || cfg.username,
    password: opts.password || process.env.KUBEPI_PASSWORD || cfg.password,
    token: opts.token || process.env.KUBEPI_TOKEN,
  };
}

const VALUE_OPT_NAMES = ['server', 'username', 'password', 'token', 'format', 'pick', 'search', 'sort', 'page', 'pageSize', 'limit', 'data'] as const;

/** A value starting with `--` means the previous flag swallowed another flag: `--search --format=json`. */
function assertNoSwallowedFlags(opts: GlobalOpts): void {
  for (const key of VALUE_OPT_NAMES) {
    const v = (opts as Record<string, unknown>)[key];
    if (typeof v === 'string' && v.startsWith('--')) {
      throw E.invalid(`option --${key.replace(/[A-Z]/g, (c) => '-' + c.toLowerCase())} requires a value (got another flag: ${v})`);
    }
  }
  for (const f of opts.filter ?? []) {
    if (f.startsWith('--')) throw E.invalid('option --filter requires a value');
  }
}

function parseErrorMessage(text: string): string {
  try {
    const json = JSON.parse(text);
    if (json && typeof json.message === 'string') return json.message;
  } catch {
    // not json
  }
  return text.trim();
}

export async function loginForToken(auth: ResolvedAuth): Promise<string> {
  if (!auth.username || !auth.password) throw E.notLoggedIn();
  let resp: Response;
  try {
    resp = await fetch(`${auth.server}${API_PREFIX}/sessions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: auth.username, password: auth.password, authMethod: 'jwt' }),
    });
  } catch (err) {
    throw E.network(errorMessage(err));
  }
  const text = await resp.text();
  if (!resp.ok) {
    const msg = parseErrorMessage(text);
    if (/mfa/i.test(msg)) throw E.mfaRequired();
    throw E.loginFailed(msg || `status ${resp.status}`);
  }
  const token = text.trim();
  // Only the token is persisted here. Account credentials are saved solely by
  // the explicit `kubepi login` action, so one-off env-var / flag usage never
  // writes a plaintext password to disk.
  patchConfig({ token, tokenSavedAt: Date.now() });
  return token;
}

export async function ensureToken(auth: ResolvedAuth, forceRefresh = false): Promise<string> {
  if (auth.token) return auth.token;
  if (!forceRefresh) {
    const cfg = loadConfig();
    if (cfg.token && cfg.tokenSavedAt && Date.now() - cfg.tokenSavedAt < TOKEN_MAX_AGE_MS) return cfg.token;
  }
  return loginForToken(auth);
}

export interface RequestOptions {
  method?: string;
  query?: Record<string, string | number | boolean | undefined>;
  body?: unknown;
  /** k8s proxy style endpoints return raw payloads instead of the {success,data} envelope. */
  raw?: boolean;
}

async function doFetch(url: string, method: string, token: string, body: unknown): Promise<Response> {
  try {
    return await fetch(url, {
      method,
      headers: {
        Authorization: `Bearer ${token}`,
        ...(body !== undefined ? { 'Content-Type': 'application/json' } : {}),
      },
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  } catch (err) {
    throw E.network(errorMessage(err));
  }
}

function buildUrl(auth: ResolvedAuth, path: string, query?: RequestOptions['query']): string {
  const url = new URL(`${auth.server}${API_PREFIX}${path}`);
  for (const [k, v] of Object.entries(query ?? {})) {
    if (v !== undefined && v !== '') url.searchParams.set(k, String(v));
  }
  return url.toString();
}

export async function request(auth: ResolvedAuth, path: string, ropts: RequestOptions = {}, retried = false): Promise<unknown> {
  const token = await ensureToken(auth);
  const method = ropts.method ?? (ropts.body !== undefined ? 'POST' : 'GET');
  const resp = await doFetch(buildUrl(auth, path, ropts.query), method, token, ropts.body);
  const text = await resp.text();

  if (resp.status === 401 && !retried && !auth.token && auth.username && auth.password) {
    await ensureToken(auth, true);
    return request(auth, path, ropts, true);
  }

  // Unknown API paths fall through to the embedded SPA and return HTML with 200.
  const contentType = resp.headers.get('content-type') ?? '';
  if (contentType.includes('text/html')) throw E.notFound(path);

  let json: unknown = text;
  try {
    json = text ? JSON.parse(text) : null;
  } catch {
    // plain text body (e.g. raw k8s responses)
  }

  if (!resp.ok) {
    const msg = json && typeof json === 'object' && 'message' in json ? String((json as { message: unknown }).message) : text.trim();
    if (resp.status === 401) throw E.authExpired();
    if (resp.status === 403) throw E.forbidden(msg);
    if (resp.status === 404) throw E.notFound(msg || path);
    if (resp.status === 400) throw E.badRequest(msg || `bad request: ${path}`);
    throw E.network(msg ? `${msg} (status ${resp.status})` : `status ${resp.status}`);
  }

  if (ropts.raw) {
    if (json && typeof json === 'object' && (json as { success?: unknown }).success === false) {
      throw E.network(String((json as { message?: unknown }).message ?? 'request failed'));
    }
    return json;
  }

  if (json && typeof json === 'object' && 'success' in json) {
    const envelope = json as { success: boolean; data?: unknown; message?: string; code?: number };
    if (!envelope.success) throw E.network(envelope.message ?? 'request failed');
    return envelope.data;
  }
  return json;
}

/** Convenience wrappers. */
export const api = {
  get: (auth: ResolvedAuth, path: string, query?: RequestOptions['query'], raw = false) =>
    request(auth, path, { method: 'GET', query, raw }),
  post: (auth: ResolvedAuth, path: string, body?: unknown, query?: RequestOptions['query'], raw = false) =>
    request(auth, path, { method: 'POST', body, query, raw }),
  put: (auth: ResolvedAuth, path: string, body?: unknown, query?: RequestOptions['query']) =>
    request(auth, path, { method: 'PUT', body, query }),
  del: (auth: ResolvedAuth, path: string) => request(auth, path, { method: 'DELETE' }),
};

/** Decode the UserProfile embedded in a KubePi JWT without verifying the signature. */
export function decodeTokenProfile(token: string): Record<string, unknown> | null {
  const parts = token.split('.');
  if (parts.length < 2) return null;
  try {
    const payload = parts[1].replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(Buffer.from(payload, 'base64').toString('utf8')) as Record<string, unknown>;
  } catch {
    return null;
  }
}
