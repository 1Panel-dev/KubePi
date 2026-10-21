import { readFileSync } from 'node:fs';
import { api } from '../client.js';
import { E } from '../errors.js';
import { getPath } from '../output.js';
import { deepMerge, registerModule, type ModuleDef, type Opts } from '../module.js';
import { collect, str } from '../shared.js';
import type { Command } from 'commander';

function readFile(path: unknown, flag: string): string {
  if (typeof path !== 'string' || path === '') throw E.missing(flag);
  try {
    return readFileSync(path, 'utf8');
  } catch {
    throw E.badRequest(`can not read file: ${path}`);
  }
}

/** Normalize auth mode aliases to KubePi values: bearer | certificate | configFile. */
function normalizeMode(v: unknown): string | undefined {
  const s = str(v);
  if (!s) return undefined;
  if (s === 'token' || s === 'bearer' || s === 'bearer-token') return 'bearer';
  if (s === 'kubeconfig' || s === 'configFile' || s === 'config-file') return 'configFile';
  if (s === 'certificate' || s === 'cert') return 'certificate';
  throw E.badRequest(`invalid --mode: ${s} (expected bearer | certificate | configFile)`);
}

function clusterCreateBody(fields: Record<string, unknown>, opts: Opts): Record<string, unknown> {
  const name = str(opts.name) ?? str(fields.name);
  if (!name) throw E.missing('--name');
  const mode =
    normalizeMode(opts.mode) ??
    normalizeMode(getPath(fields, 'spec.authentication.mode')) ??
    (opts.bearerToken || fields.bearerToken ? 'bearer' : undefined) ??
    (opts.kubeconfig ? 'configFile' : undefined) ??
    'bearer';
  const apiServer = str(opts.apiServer) ?? str(getPath(fields, 'spec.connect.forward.apiServer'));
  if (mode !== 'configFile' && !apiServer) throw E.missing('--api-server');

  const labels = (opts.label as string[] | undefined) ?? (fields.labels as string[] | undefined) ?? [];
  const body: Record<string, unknown> = {
    apiVersion: 'kubepi.io/v1',
    kind: 'Cluster',
    name,
    description: str(opts.description) ?? str(fields.description) ?? '',
    labels,
    spec: {
      connect: { direction: 'forward', forward: { apiServer: apiServer ?? '' } },
      authentication: { mode },
    },
  };
  if (mode === 'bearer') {
    const token = str(opts.bearerToken) ?? str(getPath(fields, 'spec.authentication.bearerToken'));
    if (!token) throw E.missing('--bearer-token');
    (body.spec as Record<string, unknown>).authentication = { mode, bearerToken: token };
  } else if (mode === 'configFile') {
    body.configContentStr = opts.kubeconfig ? readFile(opts.kubeconfig, '--kubeconfig') : str(fields.configContentStr);
    if (!body.configContentStr) throw E.missing('--kubeconfig or configContentStr');
  } else {
    body.certDataStr = opts.certFile ? readFile(opts.certFile, '--cert-file') : str(fields.certDataStr);
    body.keyDataStr = opts.keyFile ? readFile(opts.keyFile, '--key-file') : str(fields.keyDataStr);
    if (!body.certDataStr || !body.keyDataStr) throw E.missing('--cert-file and --key-file');
  }
  if (opts.caFile) body.caDataStr = readFile(opts.caFile, '--ca-file');
  else if (fields.caDataStr) body.caDataStr = fields.caDataStr;
  return deepMerge(fields, body) as Record<string, unknown>;
}

function clusterUpdateBody(fields: Record<string, unknown>, _current: unknown, opts: Opts): Record<string, unknown> {
  const labels = (opts.label as string[] | undefined) ?? (fields.labels as string[] | undefined);
  if (labels) return { withLabel: true, labels };
  const mode = normalizeMode(opts.mode) ?? normalizeMode(fields.mode);
  const body: Record<string, unknown> = deepMerge(fields, {}) as Record<string, unknown>;
  if (mode) body.mode = mode;
  if (opts.apiServer) body.apiServer = opts.apiServer;
  if (opts.bearerToken) body.token = opts.bearerToken;
  if (opts.kubeconfig) body.configFileContent = readFile(opts.kubeconfig, '--kubeconfig');
  if (opts.certFile) body.certData = readFile(opts.certFile, '--cert-file');
  if (opts.keyFile) body.keyData = readFile(opts.keyFile, '--key-file');
  body.withLabel = false;
  return body;
}

const clusterAuthOptions = (c: Command) => {
  c.option('--name <name>', 'cluster name');
  c.option('--mode <mode>', 'auth mode: bearer | certificate | configFile');
  c.option('--api-server <url>', 'Kubernetes apiserver URL');
  c.option('--bearer-token <token>', 'bearer token (mode=bearer)');
  c.option('--kubeconfig <path>', 'kubeconfig file path (mode=configFile)');
  c.option('--cert-file <path>', 'client certificate file (mode=certificate)');
  c.option('--key-file <path>', 'client key file (mode=certificate)');
};

/** Recursively strip sensitive keys (e.g. privateKey) from cluster payloads. */
function stripSensitive(v: unknown): unknown {
  if (Array.isArray(v)) return v.map(stripSensitive);
  if (v && typeof v === 'object') {
    const out: Record<string, unknown> = {};
    for (const [k, val] of Object.entries(v as Record<string, unknown>)) {
      if (k === 'privateKey') continue;
      if (k === 'bearerToken' && typeof val === 'string' && val.length > 0) {
        out[k] = '***';
        continue;
      }
      out[k] = stripSensitive(val);
    }
    return out;
  }
  return v;
}

export function registerCluster(program: Parameters<typeof registerModule>[0]): void {
  const cluster: ModuleDef = {
    name: 'cluster',
    description: 'manage clusters (集群)',
    basePath: '/clusters',
    searchPath: '/clusters/search',
    defaultColumns: ['name', 'status.phase', 'status.version', 'spec.connect.forward.apiServer', 'createAt'],
    sanitize: stripSensitive,
    createOptions: (c) => {
      clusterAuthOptions(c);
      c.option('--description <text>', 'description');
      c.option('--ca-file <path>', 'CA certificate file');
      c.option('--label <label>', 'label, repeatable', collect);
    },
    updateOptions: (c) => {
      clusterAuthOptions(c);
      c.option('--label <label>', 'replace labels, repeatable', collect);
    },
    buildCreate: clusterCreateBody,
    buildUpdate: clusterUpdateBody,
    help: [
      'Examples:',
      '  kubepi cluster create --name=dev --mode=bearer --api-server=https://1.2.3.4:6443 --bearer-token=<token>',
      '  kubepi cluster create --name=dev --mode=configFile --kubeconfig=~/.kube/config',
      '  kubepi cluster update dev --label=env=dev --label=team=infra',
    ].join('\n'),
    actions: [
      {
        name: 'namespaces <name>',
        description: 'list namespaces of a cluster',
        run: async (args, opts, auth) => api.get(auth, `/clusters/${encodeURIComponent(String(args[0]))}/namespaces`),
      },
      {
        name: 'members <name>',
        description: 'list members of a cluster (see also: clustermember)',
        run: async (args, opts, auth) => api.get(auth, `/clusters/${encodeURIComponent(String(args[0]))}/members`),
        printColumns: ['name', 'clusterRoles', 'createAt'],
      },
      {
        name: 'repos <name>',
        description: 'list chart repos bound to a cluster (see also: clusterrepo)',
        options: (c) => c.option('--detail', 'include registry details'),
        run: async (args, opts, auth) =>
          api.get(auth, `/clusters/${encodeURIComponent(String(args[0]))}/repos${opts.detail ? '/detail' : ''}`),
      },
    ],
  };
  registerModule(program, cluster);

  registerModule(program, {
    name: 'clustermember',
    description: 'manage cluster members (集群成员), requires --cluster',
    basePath: (opts) => `/clusters/${encodeURIComponent(String(opts.cluster))}/members`,
    requires: ['cluster'],
    moduleOptions: (c) => c.option('--cluster <name>', 'cluster name'),
    defaultColumns: ['name', 'clusterRoles', 'createAt'],
    createOptions: (c) => {
      c.option('--name <name>', 'username');
      c.option('--cluster-role <role>', 'cluster role, repeatable', collect);
      c.option('--namespace-role <ns:role1,role2>', 'namespace roles, repeatable', collect);
    },
    updateOptions: (c) => {
      c.option('--cluster-role <role>', 'cluster role, repeatable', collect);
      c.option('--namespace-role <ns:role1,role2>', 'namespace roles, repeatable', collect);
    },
    buildCreate: (fields, opts) => memberBody(fields, opts),
    buildUpdate: (fields, current, opts) => deepMerge(current ?? {}, memberBody(fields, opts, true)),
    help: [
      'Examples:',
      '  kubepi clustermember --cluster=dev',
      '  kubepi clustermember create --cluster=dev --name=tom --cluster-role=cluster-admin',
      '  kubepi clustermember create --cluster=dev --name=tom --namespace-role=default:view,edit',
    ].join('\n'),
  });

  registerModule(program, {
    name: 'clusterrole',
    description: 'manage K8s ClusterRoles inside a cluster, requires --cluster',
    basePath: (opts) => `/clusters/${encodeURIComponent(String(opts.cluster))}/clusterroles`,
    requires: ['cluster'],
    moduleOptions: (c) => c.option('--cluster <name>', 'cluster name'),
    defaultColumns: ['metadata.name', 'metadata.creationTimestamp'],
    noGet: true,
    noPrefill: true,
    createOptions: (c) => {
      c.option('--name <name>', 'ClusterRole name');
      c.option('--rules <json>', 'rules JSON array, e.g. [{"apiGroups":[""],"resources":["pods"],"verbs":["get"]}]');
    },
    updateOptions: (c) => {
      c.option('--rules <json>', 'rules JSON array');
    },
    buildCreate: (fields, opts) => clusterRoleBody(fields, opts),
    buildUpdate: (fields, _current, opts) => clusterRoleBody(fields, opts),
    help: [
      'Examples:',
      `  kubepi clusterrole create --cluster=dev --name=pod-reader --rules='[{"apiGroups":[""],"resources":["pods"],"verbs":["get","list"]}]'`,
    ].join('\n'),
  });

  registerModule(program, {
    name: 'clusterrepo',
    description: 'manage chart repos bound to a cluster, requires --cluster',
    basePath: (opts) => `/clusters/${encodeURIComponent(String(opts.cluster))}/repos`,
    requires: ['cluster'],
    moduleOptions: (c) => c.option('--cluster <name>', 'cluster name'),
    defaultColumns: ['repo', 'cluster', 'createAt'],
    noGet: true,
    noUpdate: true,
    createOptions: (c) => {
      c.option('--repo <repo>', 'repo name (from chartrepo), repeatable', collect);
    },
    buildCreate: (fields, opts) => {
      const repos = (opts.repo as string[] | undefined) ?? (fields.repos as string[] | undefined);
      if (!repos || repos.length === 0) throw E.missing('--repo');
      return { repos, cluster: opts.cluster };
    },
    actions: [
      {
        name: 'detail',
        description: 'list repos with registry details',
        run: async (_args, opts, auth) => api.get(auth, `/clusters/${encodeURIComponent(String(opts.cluster))}/repos/detail`),
      },
    ],
  });
}

function memberBody(fields: Record<string, unknown>, opts: Opts, partial = false): Record<string, unknown> {
  const body: Record<string, unknown> = { ...fields };
  const name = str(opts.name) ?? str(fields.name);
  if (!partial && !name) throw E.missing('--name');
  if (name) body.name = name;
  const clusterRoles = (opts.clusterRole as string[] | undefined) ?? (fields.clusterRoles as string[] | undefined);
  if (clusterRoles) body.clusterRoles = clusterRoles;
  const nsRoles = (opts.namespaceRole as string[] | undefined) ?? (fields.namespaceRoles as { namespace: string; roles: string[] }[] | undefined);
  if (nsRoles && nsRoles.length > 0 && typeof nsRoles[0] === 'string') {
    body.namespaceRoles = (nsRoles as unknown as string[]).map((item) => {
      const [namespace, roles] = item.split(':');
      if (!namespace || !roles) throw E.badRequest(`invalid --namespace-role: ${item} (expected ns:role1,role2)`);
      return { namespace, roles: roles.split(',').filter(Boolean) };
    });
  } else if (nsRoles) {
    body.namespaceRoles = nsRoles;
  }
  if (!partial && !body.clusterRoles && !body.namespaceRoles) throw E.missing('--cluster-role or --namespace-role');
  return body;
}

function clusterRoleBody(fields: Record<string, unknown>, opts: Opts): Record<string, unknown> {
  if (fields.metadata || fields.rules) {
    if (opts.name) setNameOnBody(fields, opts);
    return fields;
  }
  const name = str(opts.name);
  if (!name) throw E.missing('--name');
  let rules: unknown;
  try {
    rules = JSON.parse(String(opts.rules ?? ''));
  } catch {
    throw E.badRequest('--rules is not valid JSON');
  }
  if (!Array.isArray(rules) || rules.length === 0) throw E.missing('--rules');
  return { apiVersion: 'rbac.authorization.k8s.io/v1', kind: 'ClusterRole', metadata: { name }, rules };
}

function setNameOnBody(fields: Record<string, unknown>, opts: Opts): void {
  if (typeof fields.metadata === 'object' && fields.metadata !== null) {
    (fields.metadata as Record<string, unknown>).name = opts.name;
  } else {
    fields.metadata = { name: opts.name };
  }
}
