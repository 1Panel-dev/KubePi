import { api } from '../client.js';
import { E } from '../errors.js';
import { collectFields, deepMerge, registerModule } from '../module.js';
import { parseFilters, buildConditions } from '../query.js';
import { str } from '../shared.js';

async function ssoCurrent(auth: Parameters<typeof api.get>[0]): Promise<Record<string, unknown> | undefined> {
  let res: unknown;
  try {
    res = await api.get(auth, '/sso/');
  } catch (err) {
    // On a fresh install the storm bucket does not exist yet and the server
    // answers 500 "bucket not found": treat as "no config yet" (→ POST).
    if (err instanceof Error && err.message.includes('not found')) return undefined;
    throw err;
  }
  if (Array.isArray(res) && res.length > 0) return res[0] as Record<string, unknown>;
  if (res && typeof res === 'object') return res as Record<string, unknown>;
  return undefined;
}

export function registerSystem(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'sso',
    description: 'manage SSO settings (single config)',
    basePath: '/sso/',
    noList: true,
    noGet: true,
    noCreate: true,
    noUpdate: true,
    noDelete: true,
    actions: [
      {
        name: 'get',
        description: 'show the SSO config',
        run: async (_args, _opts, auth) => api.get(auth, '/sso/'),
      },
      {
        name: 'set',
        description: 'create or update the SSO config (PUT when one exists, POST otherwise)',
        run: async (_args, opts, auth, rawArgs) => {
          const body = collectFields(opts, rawArgs);
          if (Object.keys(body).length === 0) throw E.missing('--data or --field=value');
          const current = await ssoCurrent(auth);
          // The endpoint replaces the whole config: merge current values so
          // fields not given on the command line are preserved.
          return current ? api.put(auth, '/sso/', deepMerge(current, body)) : api.post(auth, '/sso/', body);
        },
      },
      {
        name: 'status',
        description: 'show SSO login status',
        run: async (_args, _opts, auth) => api.get(auth, '/sso/status'),
      },
    ],
    help: [
      'Examples:',
      `  kubepi sso set --data='{"protocol":"oidc","interfaceAddress":"http://kubepi.local","authServerAddress":"http://oidc.local","clientId":"kubepi","clientSecret":"xxx"}'`,
    ].join('\n'),
  });

  const logQuery = async (path: string, opts: Record<string, unknown> & { filter?: string[]; search?: string; page?: string; pageSize?: string }, auth: Parameters<typeof api.post>[0]) => {
    const groups = parseFilters(opts.filter);
    const conditions = buildConditions(groups, str(opts.search)) ?? {};
    const query = opts.page || opts.pageSize ? { pageNum: opts.page ?? '1', pageSize: opts.pageSize ?? '10' } : undefined;
    const res = await api.post(auth, path, { conditions }, query);
    if (res && typeof res === 'object' && !Array.isArray(res) && 'total' in res) {
      return (res as { items?: unknown[] }).items ?? [];
    }
    return res;
  };

  registerModule(program, {
    name: 'log',
    description: 'query login / operation logs',
    noList: true,
    noGet: true,
    noCreate: true,
    noUpdate: true,
    noDelete: true,
    actions: [
      {
        name: 'login',
        description: 'query login logs (supports --search/--filter/--page)',
        printColumns: ['userName', 'ip', 'city', 'createAt'],
        run: async (_args, opts, auth) => logQuery('/systems/login/logs/search', opts, auth),
      },
      {
        name: 'operation',
        description: 'query operation logs (supports --search/--filter/--page)',
        printColumns: ['operator', 'operation', 'operationDomain', 'specificInformation', 'createAt'],
        run: async (_args, opts, auth) => logQuery('/systems/operation/logs/search', opts, auth),
      },
    ],
    help: ['Examples:', '  kubepi log login --search=admin', "  kubepi log operation --filter='operator:admin' --page=1 --page-size=20"].join('\n'),
  });
}
