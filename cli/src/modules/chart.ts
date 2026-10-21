import type { Command } from 'commander';
import { api } from '../client.js';
import { E } from '../errors.js';
import { deepMerge, registerModule } from '../module.js';
import { setPath } from '../output.js';
import { collect, str } from '../shared.js';

function clusterOpt(c: Command): void {
  c.option('--cluster <name>', 'cluster name');
}

export function registerChart(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'chartrepo',
    description: 'manage helm chart repos (per cluster), requires --cluster',
    basePath: (opts) => `/charts/${encodeURIComponent(String(opts.cluster))}/repos`,
    requires: ['cluster'],
    moduleOptions: clusterOpt,
    defaultColumns: ['name', 'url'],
    createOptions: (c) => {
      c.option('--name <name>', 'repo name');
      c.option('--url <url>', 'repo index URL');
      c.option('--repo-username <name>', 'repo username');
      c.option('--repo-password <password>', 'repo password');
    },
    updateOptions: (c) => {
      c.option('--url <url>', 'repo index URL');
      c.option('--repo-username <name>', 'repo username');
      c.option('--repo-password <password>', 'repo password');
    },
    buildCreate: (fields, opts) => {
      const name = str(opts.name) ?? str(fields.name);
      if (!name) throw E.missing('--name');
      const url = str(opts.url) ?? str(fields.url);
      if (!url) throw E.missing('--url');
      return deepMerge(fields, {
        name,
        url,
        userName: str(opts.repoUsername) ?? str(fields.userName) ?? '',
        password: str(opts.repoPassword) ?? str(fields.password) ?? '',
      });
    },
    buildUpdate: (fields, current, opts) => {
      const patch = deepMerge(current ?? {}, fields) as Record<string, unknown>;
      if (opts.url) patch.url = opts.url;
      if (opts.repoUsername) patch.userName = opts.repoUsername;
      if (opts.repoPassword) patch.password = opts.repoPassword;
      return patch;
    },
    actions: [
      {
        name: 'sync <name>',
        description: 'sync a repo (fetch latest index)',
        run: async (args, opts, auth) =>
          api.post(auth, `/charts/${encodeURIComponent(String(opts.cluster))}/repos/sync/${encodeURIComponent(String(args[0]))}`),
      },
    ],
    help: ['Examples:', '  kubepi chartrepo create --cluster=dev --name=bitnami --url=https://charts.bitnami.com/bitnami'].join('\n'),
  });

  registerModule(program, {
    name: 'chart',
    description: 'browse and install helm charts (per cluster), requires --cluster',
    basePath: (opts) => `/charts/${encodeURIComponent(String(opts.cluster))}/search`,
    requires: ['cluster'],
    moduleOptions: (c) => {
      clusterOpt(c);
      c.option('--repo <name>', 'filter by repo');
    },
    listQuery: (opts) => ({ pattern: str(opts.search), repo: str(opts.repo) }),
    defaultColumns: ['name', 'repo', 'description'],
    noGet: true,
    noCreate: true,
    noUpdate: true,
    noDelete: true,
    actions: [
      {
        name: 'detail <name>',
        description: 'show chart versions/values; --version picks one version',
        options: (c) => {
          c.option('--repo <name>', 'repo name');
          c.option('--version <version>', 'chart version');
        },
        run: async (args, opts, auth) => {
          const cluster = encodeURIComponent(String(opts.cluster));
          const name = encodeURIComponent(String(args[0]));
          if (opts.version) return api.get(auth, `/charts/${cluster}/detail/${name}`, { repo: str(opts.repo), version: str(opts.version) });
          return api.get(auth, `/charts/${cluster}/${name}`, { repo: str(opts.repo) });
        },
      },
      {
        name: 'install',
        description: 'install a chart release',
        options: (c) => {
          c.option('--name <name>', 'release name');
          c.option('--chart-name <name>', 'chart name');
          c.option('--repo <name>', 'repo name');
          c.option('--chart-version <version>', 'chart version');
          c.option('--namespace <ns>', 'target namespace');
          c.option('--values <json>', 'values JSON object');
          c.option('--set <key=value>', 'set a values key, repeatable', collect);
        },
        run: async (_args, opts, auth) => {
          if (!str(opts.name)) throw E.missing('--name');
          if (!str(opts.chartName)) throw E.missing('--chart-name');
          if (!str(opts.repo)) throw E.missing('--repo');
          if (!str(opts.chartVersion)) throw E.missing('--chart-version');
          let values: Record<string, unknown> = {};
          if (opts.values) {
            try {
              values = JSON.parse(String(opts.values));
            } catch {
              throw E.badRequest('--values is not valid JSON');
            }
          }
          for (const pair of (opts.set as string[] | undefined) ?? []) {
            const eq = pair.indexOf('=');
            if (eq <= 0) throw E.badRequest(`invalid --set: ${pair}`);
            let v: unknown = pair.slice(eq + 1);
            try {
              v = JSON.parse(String(v));
            } catch {
              // keep string
            }
            setPath(values, pair.slice(0, eq), v);
          }
          const body = {
            name: opts.name,
            chartName: opts.chartName,
            repo: opts.repo,
            chartVersion: opts.chartVersion,
            cluster: opts.cluster,
            namespace: str(opts.namespace) ?? 'default',
            values,
          };
          return api.post(auth, `/charts/${encodeURIComponent(String(opts.cluster))}/install`, body);
        },
      },
      {
        name: 'installed',
        description: 'list installed releases (apps) in the cluster',
        run: async (_args, opts, auth) => {
          const query = opts.page || opts.pageSize ? { pageNum: str(opts.page) ?? '1', pageSize: str(opts.pageSize) ?? '10' } : undefined;
          const res = await api.get(auth, `/apps/${encodeURIComponent(String(opts.cluster))}/search`, query);
          if (res && typeof res === 'object' && !Array.isArray(res) && Array.isArray((res as { items?: unknown[] }).items)) {
            return (res as { items: unknown[] }).items;
          }
          return res;
        },
        printColumns: ['name', 'namespace', 'chart.metadata.name', 'chart.metadata.version', 'info.status'],
      },
      {
        name: 'uninstall <name>',
        description: 'uninstall a release',
        options: (c) => c.option('--namespace <ns>', 'release namespace'),
        run: async (args, opts, auth) =>
          api.del(
            auth,
            `/apps/${encodeURIComponent(String(opts.cluster))}/${encodeURIComponent(str(opts.namespace) ?? 'default')}/${encodeURIComponent(String(args[0]))}`,
          ),
      },
    ],
    help: [
      'Examples:',
      '  kubepi chart --cluster=dev --repo=bitnami --search=nginx',
      '  kubepi chart detail nginx --cluster=dev --repo=bitnami --version=15.0.0',
      '  kubepi chart install --cluster=dev --name=my-nginx --chart-name=nginx --repo=bitnami --chart-version=15.0.0 --namespace=default',
    ].join('\n'),
  });
}
