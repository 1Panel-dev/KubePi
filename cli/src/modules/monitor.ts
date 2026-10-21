import { api } from '../client.js';
import { E } from '../errors.js';
import { collectFields, deepMerge, registerModule } from '../module.js';
import { buildConditions, parseFilters } from '../query.js';
import { boolOpt, str } from '../shared.js';

async function grafanaCurrent(auth: Parameters<typeof api.get>[0]): Promise<Record<string, unknown> | undefined> {
  let res: unknown;
  try {
    res = await api.get(auth, '/monitor/grafana/');
  } catch (err) {
    // On a fresh install the storm bucket does not exist yet and the server
    // answers 500 "bucket not found": treat as "no config yet" (→ POST).
    if (err instanceof Error && err.message.includes('not found')) return undefined;
    throw err;
  }
  if (Array.isArray(res) && res.length > 0) return res[0] as Record<string, unknown>;
  if (res && typeof res === 'object' && !Array.isArray(res)) return res as Record<string, unknown>;
  return undefined;
}

function metricsBody(fields: Record<string, unknown>, opts: Record<string, unknown>): Record<string, unknown> {
  const body = deepMerge(fields, {}) as Record<string, unknown>;
  if (str(opts.name)) body.name = opts.name;
  if (str(opts.type)) body.type = opts.type;
  if (str(opts.endpoint)) body.endPoint = opts.endpoint;
  const auth = boolOpt(opts.auth);
  if (auth !== undefined) body.auth = auth;
  const username = str(opts.repoUsername);
  const password = str(opts.repoPassword);
  if (username !== undefined || password !== undefined) {
    const cur = (body.credential ?? {}) as Record<string, unknown>;
    body.credential = { username: username ?? cur.username ?? '', password: password ?? cur.password ?? '' };
  }
  return body;
}

const metricOptions = (c: { option: (flags: string, desc?: string) => unknown }) => {
  c.option('--name <name>', 'metrics source name');
  c.option('--type <type>', 'type, e.g. prometheus');
  c.option('--endpoint <url>', 'endpoint URL');
  c.option('--repo-username <name>', 'auth username');
  c.option('--repo-password <password>', 'auth password');
  c.option('--auth <bool>', 'requires auth, true/false');
};

export function registerMonitor(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'monitor',
    description: 'manage grafana / metrics sources (requires KubePi >= v2.0)',
    noList: true,
    noGet: true,
    noCreate: true,
    noUpdate: true,
    noDelete: true,
    actions: [
      {
        name: 'grafana',
        description: 'show the grafana config',
        run: async (_args, _opts, auth) => api.get(auth, '/monitor/grafana/'),
      },
      {
        name: 'grafana-set',
        description: 'create or update the grafana config',
        options: (c) => {
          c.option('--address <url>', 'grafana URL');
          c.option('--enable <bool>', 'enable, true/false');
          c.option('--default-dashboard <bool>', 'import default dashboards, true/false');
          c.option('--sa-token <token>', 'grafana service account token');
        },
        run: async (_args, opts, auth, rawArgs) => {
          const fields = collectFields(opts, rawArgs);
          const body: Record<string, unknown> = { ...fields };
          if (str(opts.address)) body.address = opts.address;
          const enable = boolOpt(opts.enable);
          if (enable !== undefined) body.enable = enable;
          const dd = boolOpt(opts.defaultDashboard);
          if (dd !== undefined) body.default_dashboard = dd;
          if (str(opts.saToken)) body.service_account_token = opts.saToken;
          if (Object.keys(body).length === 0) throw E.missing('--address/--enable/--sa-token or --data');
          const current = await grafanaCurrent(auth);
          // The endpoint replaces the whole config: merge current values so
          // fields not given on the command line are preserved.
          return current ? api.put(auth, '/monitor/grafana/', deepMerge(current, body)) : api.post(auth, '/monitor/grafana/', body);
        },
      },
      {
        name: 'grafana-test',
        description: 'test grafana connectivity with given (or stored) config',
        run: async (_args, opts, auth, rawArgs) => {
          const fields = collectFields(opts, rawArgs);
          const body = Object.keys(fields).length > 0 ? fields : ((await grafanaCurrent(auth)) ?? ({}));
          return api.post(auth, '/monitor/grafana/test/connect', body);
        },
      },
      {
        name: 'grafana-import',
        description: 'import default dashboards into grafana',
        run: async (_args, opts, auth, rawArgs) => {
          const fields = collectFields(opts, rawArgs);
          const body = Object.keys(fields).length > 0 ? fields : ((await grafanaCurrent(auth)) ?? ({}));
          return api.post(auth, '/monitor/grafana/import', body);
        },
      },
      {
        name: 'metrics',
        description: 'list metrics sources (supports --search/--filter/--page)',
        printColumns: ['name', 'type', 'endPoint', 'auth', 'createAt'],
        run: async (_args, opts, auth, rawArgs) => {
          const groups = parseFilters(opts.filter);
          const conditions = buildConditions(groups, str(opts.search)) ?? {};
          const query = opts.page || opts.pageSize ? { pageNum: opts.page ?? '1', pageSize: opts.pageSize ?? '10' } : undefined;
          const res = await api.post(auth, '/monitor/metrics/search', { conditions }, query);
          if (res && typeof res === 'object' && !Array.isArray(res) && 'total' in res) {
            return (res as { items?: unknown[] }).items ?? [];
          }
          return res;
        },
      },
      {
        name: 'metrics-get <name>',
        description: 'show a metrics source',
        run: async (args, _opts, auth) => api.get(auth, `/monitor/metrics/${encodeURIComponent(String(args[0]))}`),
      },
      {
        name: 'metrics-create',
        description: 'create a metrics source',
        options: metricOptions,
        run: async (_args, opts, auth, rawArgs) => {
          const body = metricsBody(collectFields(opts, rawArgs), opts);
          if (!str(body.name)) throw E.missing('--name');
          if (!str(body.endPoint)) throw E.missing('--endpoint');
          body.type = str(body.type) ?? 'prometheus';
          return api.post(auth, '/monitor/metrics/', body);
        },
      },
      {
        name: 'metrics-update <name>',
        description: 'update a metrics source',
        options: metricOptions,
        run: async (args, opts, auth, rawArgs) => {
          const current = (await api.get(auth, `/monitor/metrics/${encodeURIComponent(String(args[0]))}`)) as Record<string, unknown>;
          const body = metricsBody(deepMerge(current, collectFields(opts, rawArgs)) as Record<string, unknown>, opts);
          return api.put(auth, `/monitor/metrics/${encodeURIComponent(String(args[0]))}`, body);
        },
      },
      {
        name: 'metrics-delete <name>',
        description: 'delete a metrics source',
        run: async (args, _opts, auth) => api.del(auth, `/monitor/metrics/${encodeURIComponent(String(args[0]))}`),
      },
      {
        name: 'metrics-query <name>',
        description: 'run a PromQL instant query against a metrics source',
        options: (c) => {
          c.option('--promql <expr>', 'PromQL expression');
          c.option('--time <rfc3339>', 'evaluation timestamp');
        },
        run: async (args, opts, auth, rawArgs) => {
          if (!str(opts.promql)) throw E.missing('--promql');
          return api.get(auth, `/monitor/metrics/${encodeURIComponent(String(args[0]))}/query`, {
            promql: String(opts.promql),
            time: str(opts.time),
          });
        },
      },
    ],
    help: [
      'Examples:',
      '  kubepi monitor grafana-set --address=http://grafana.local --sa-token=xxx --enable=true --default-dashboard=true',
      "  kubepi monitor metrics-query prom --promql='up'",
    ].join('\n'),
  });
}
