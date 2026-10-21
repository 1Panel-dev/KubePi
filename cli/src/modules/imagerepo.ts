import { api } from '../client.js';
import { E } from '../errors.js';
import { deepMerge, registerModule } from '../module.js';
import { boolOpt, numOpt, str } from '../shared.js';

export function registerImageRepo(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'imagerepo',
    description: 'manage image registries (镜像仓库)',
    basePath: '/imagerepos',
    searchPath: '/imagerepos/search',
    defaultColumns: ['name', 'type', 'endPoint', 'repoName', 'auth', 'createAt'],
    createOptions: (c) => {
      c.option('--name <name>', 'registry name');
      c.option('--type <type>', 'registry type, e.g. Harbor');
      c.option('--endpoint <url>', 'registry endpoint URL');
      c.option('--download-url <url>', 'image pull URL presented to users');
      c.option('--repo-name <name>', 'repository/project name inside the registry');
      c.option('--repo-username <name>', 'registry username');
      c.option('--repo-password <password>', 'registry password');
      c.option('--auth <bool>', 'registry requires auth, true/false');
      c.option('--allow-anonymous <bool>', 'allow anonymous pull, true/false');
      c.option('--version <version>', 'registry version');
    },
    updateOptions: (c) => {
      c.option('--type <type>', 'registry type');
      c.option('--endpoint <url>', 'registry endpoint URL');
      c.option('--download-url <url>', 'image pull URL');
      c.option('--repo-name <name>', 'repository/project name');
      c.option('--repo-username <name>', 'registry username');
      c.option('--repo-password <password>', 'registry password');
      c.option('--auth <bool>', 'true/false');
      c.option('--allow-anonymous <bool>', 'true/false');
      c.option('--version <version>', 'registry version');
    },
    buildCreate: (fields, opts) => {
      const name = str(opts.name) ?? str(fields.name);
      if (!name) throw E.missing('--name');
      const endPoint = str(opts.endpoint) ?? str(fields.endPoint);
      if (!endPoint) throw E.missing('--endpoint');
      const username = str(opts.repoUsername);
      const password = str(opts.repoPassword);
      return deepMerge(fields, {
        apiVersion: 'v1',
        kind: 'ImageRepo',
        name,
        type: str(opts.type) ?? str(fields.type) ?? 'Harbor',
        endPoint,
        downloadUrl: str(opts.downloadUrl) ?? str(fields.downloadUrl) ?? endPoint,
        repoName: str(opts.repoName) ?? str(fields.repoName) ?? '',
        auth: boolOpt(opts.auth) ?? (username ? true : false),
        allowAnonymous: boolOpt(opts.allowAnonymous) ?? false,
        version: str(opts.version) ?? str(fields.version) ?? '',
        credential: { username: username ?? '', password: password ?? '' },
      });
    },
    buildUpdate: (fields, current, opts) => {
      const patch = deepMerge(current ?? {}, fields) as Record<string, unknown>;
      if (opts.type) patch.type = opts.type;
      if (opts.endpoint) patch.endPoint = opts.endpoint;
      if (opts.downloadUrl) patch.downloadUrl = opts.downloadUrl;
      if (opts.repoName) patch.repoName = opts.repoName;
      if (opts.version) patch.version = opts.version;
      const auth = boolOpt(opts.auth);
      if (auth !== undefined) patch.auth = auth;
      const allowAnonymous = boolOpt(opts.allowAnonymous);
      if (allowAnonymous !== undefined) patch.allowAnonymous = allowAnonymous;
      const username = str(opts.repoUsername);
      const password = str(opts.repoPassword);
      if (username !== undefined || password !== undefined) {
        const cur = (patch.credential ?? {}) as Record<string, unknown>;
        patch.credential = { username: username ?? cur.username ?? '', password: password ?? cur.password ?? '' };
      }
      return patch;
    },
    actions: [
      {
        name: 'images <name>',
        description: 'list images of a registry cached in a cluster',
        options: (c) => c.option('--cluster <name>', 'cluster name').option('--search <keyword>', 'keyword').option('--page <n>').option('--limit <n>'),
        run: async (args, opts, auth) => {
          if (!opts.cluster) throw E.missing('--cluster');
          return api.get(
            auth,
            `/imagerepos/images/${encodeURIComponent(String(opts.cluster))}/${encodeURIComponent(String(args[0]))}`,
          );
        },
      },
      {
        name: 'search-images <name>',
        description: 'search images inside the registry itself (paged)',
        options: (c) => c.option('--search <keyword>', 'keyword').option('--page <n>').option('--limit <n>').option('--continue-token <token>'),
        run: async (args, opts, auth) => {
          const repo = (await api.get(auth, `/imagerepos/${encodeURIComponent(String(args[0]))}`)) as Record<string, unknown>;
          const body = {
            ...repo,
            page: numOpt(opts.page, '--page') ?? 1,
            limit: numOpt(opts.limit, '--limit') ?? 10,
            search: str(opts.search) ?? '',
            continueToken: str(opts.continueToken) ?? '',
          };
          return api.post(auth, `/imagerepos/images/${encodeURIComponent(String(args[0]))}/search`, body);
        },
      },
      {
        name: 'repositories <name>',
        description: 'list repositories/projects inside the registry',
        options: (c) => c.option('--search <keyword>', 'keyword').option('--page <n>').option('--limit <n>'),
        run: async (args, opts, auth) => {
          const repo = (await api.get(auth, `/imagerepos/${encodeURIComponent(String(args[0]))}`)) as Record<string, unknown>;
          const body = {
            ...repo,
            page: numOpt(opts.page, '--page') ?? 1,
            limit: numOpt(opts.limit, '--limit') ?? 10,
            search: str(opts.search) ?? '',
          };
          return api.post(auth, '/imagerepos/repositories/search', body);
        },
      },
    ],
    help: [
      'Examples:',
      '  kubepi imagerepo create --name=harbor --type=Harbor --endpoint=http://harbor.local --repo-username=admin --repo-password=xxx',
      '  kubepi imagerepo repositories harbor --search=library',
    ].join('\n'),
  });
}
