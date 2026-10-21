import { readFileSync } from 'node:fs';
import { api } from '../client.js';
import { E } from '../errors.js';
import { collectFields, deepMerge, registerModule, type Opts } from '../module.js';
import { boolOpt, numOpt, str } from '../shared.js';

function ldapFields(fields: Record<string, unknown>, opts: Opts): Record<string, unknown> {
  const body = deepMerge(fields, {}) as Record<string, unknown>;
  const mapping: Record<string, unknown> = {
    address: opts.address,
    port: opts.port !== undefined ? String(opts.port) : undefined,
    dn: opts.dn,
    filter: opts.filter,
    mapping: opts.mapping,
    username: opts.bindUsername,
    password: opts.bindPassword,
    tls: boolOpt(opts.tls),
    enable: boolOpt(opts.enable),
    sizeLimit: numOpt(opts.sizeLimit, '--size-limit'),
    timeLimit: numOpt(opts.timeLimit, '--time-limit'),
  };
  for (const [k, v] of Object.entries(mapping)) if (v !== undefined) body[k] = v;
  return body;
}

async function currentLdap(auth: Parameters<typeof api.get>[0]): Promise<Record<string, unknown> | undefined> {
  const list = (await api.get(auth, '/ldap/')) as unknown[];
  return Array.isArray(list) && list.length > 0 ? (list[0] as Record<string, unknown>) : undefined;
}

export function registerLdap(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'ldap',
    description: 'manage LDAP settings (single config)',
    basePath: '/ldap/',
    noGet: true,
    noUpdate: true,
    noDelete: true,
    defaultColumns: ['address', 'port', 'dn', 'enable', 'tls'],
    createOptions: (c) => {
      c.option('--address <host>', 'LDAP server address');
      c.option('--port <port>', 'LDAP port');
      c.option('--dn <dn>', 'base DN, e.g. dc=example,dc=com');
      c.option('--filter <filter>', 'user filter, e.g. (uid=%s)');
      c.option('--mapping <json>', 'attribute mapping JSON');
      c.option('--bind-username <name>', 'bind username');
      c.option('--bind-password <password>', 'bind password');
      c.option('--tls <bool>', 'use TLS, true/false');
      c.option('--enable <bool>', 'enable LDAP login, true/false');
      c.option('--size-limit <n>', 'search size limit');
      c.option('--time-limit <n>', 'search time limit (seconds)');
    },
    buildCreate: (fields, opts) => {
      const body = ldapFields(fields, opts);
      for (const key of ['address', 'port', 'dn']) {
        if (!str(body[key])) throw E.missing(`--${key}`);
      }
      return body;
    },
    actions: [
      {
        name: 'update',
        description: 'update the LDAP config (full replace via PUT)',
        options: (c) => {
          c.option('--address <host>', 'LDAP server address');
          c.option('--port <port>', 'LDAP port');
          c.option('--dn <dn>', 'base DN');
          c.option('--filter <filter>', 'user filter');
          c.option('--mapping <json>', 'attribute mapping JSON');
          c.option('--bind-username <name>', 'bind username');
          c.option('--bind-password <password>', 'bind password');
          c.option('--tls <bool>', 'true/false');
          c.option('--enable <bool>', 'true/false');
          c.option('--size-limit <n>', 'search size limit');
          c.option('--time-limit <n>', 'search time limit');
        },
        run: async (_args, opts, auth, rawArgs) => {
          const current = (await currentLdap(auth)) ?? {};
          return api.put(auth, '/ldap/', ldapFields(deepMerge(current, collectFields(opts, rawArgs)) as Record<string, unknown>, opts));
        },
      },
      {
        name: 'sync',
        description: 'sync LDAP users into KubePi',
        run: async (_args, _opts, auth) => api.post(auth, '/ldap/sync'),
      },
      {
        name: 'test-connect',
        description: 'test connection with given (or stored) LDAP settings',
        run: async (_args, opts, auth, rawArgs) => {
          const fields = collectFields(opts, rawArgs);
          const body = Object.keys(fields).length > 0 ? fields : ((await currentLdap(auth)) ?? ({}));
          return api.post(auth, '/ldap/test/connect', body);
        },
      },
      {
        name: 'test-login',
        description: 'test an LDAP user login',
        options: (c) => {
          c.option('--test-username <name>', 'LDAP username');
          c.option('--test-password <password>', 'LDAP password');
        },
        run: async (_args, opts, auth, rawArgs) => {
          const fields = collectFields(opts, rawArgs);
          const username = str(opts.testUsername) ?? str(fields.username);
          const password = str(opts.testPassword) ?? str(fields.password);
          if (!username || !password) throw E.missing('--test-username and --test-password');
          return api.post(auth, '/ldap/test/login', { username, password });
        },
      },
      {
        name: 'import',
        description: 'import selected LDAP users (--file users.json or --data)',
        options: (c) => c.option('--file <path>', 'JSON file: [{"name","email","nickName"}]'),
        run: async (_args, opts, auth, rawArgs) => {
          let users: unknown;
          if (opts.file) users = JSON.parse(readFileSync(String(opts.file), 'utf8'));
          else {
            const fields = collectFields(opts, rawArgs);
            users = fields.users;
          }
          if (!Array.isArray(users) || users.length === 0) throw E.missing('--file or --data with users array');
          return api.post(auth, '/ldap/import', { users });
        },
      },
    ],
    help: [
      'Examples:',
      '  kubepi ldap create --address=ldap.local --port=389 --dn=dc=example,dc=com --filter=(uid=%s) --bind-username=admin --bind-password=xxx --enable=true',
      '  kubepi ldap test-login --test-username=tom --test-password=xxx',
    ].join('\n'),
  });
}
