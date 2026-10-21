import { Command } from 'commander';
import { request, resolveAuth, type GlobalOpts } from '../client.js';
import { E } from '../errors.js';
import { confirmAction, mergedOpts } from '../module.js';
import { printData } from '../output.js';
import { matchGroups, matchSearch, parseFilters, applySort } from '../query.js';
import { numOpt } from '../shared.js';

interface K8sOpts extends GlobalOpts {
  namespace?: string;
  method?: string;
}

/**
 * Generic Kubernetes API passthrough:
 *   kubepi k8s <cluster> <apiPath> [--namespace=] [--method=GET] [--data=JSON]
 * Examples:
 *   kubepi k8s dev /api/v1/pods --namespace=default
 *   kubepi k8s dev /apis/apps/v1/deployments --namespace=default --search=nginx
 *   kubepi k8s dev /api/v1/namespaces --method=POST --data='{"metadata":{"name":"demo"}}'
 */
export function registerK8s(program: Command): void {
  const c = program
    .command('k8s <cluster> <path>')
    .description('call any Kubernetes API of a cluster through the KubePi proxy')
    .option('--namespace <ns>', 'namespace for namespaced resources')
    .option('--method <method>', 'HTTP method: GET (default) | POST | PUT | PATCH | DELETE')
    .addHelpText(
      'after',
      [
        'Notes:',
        '  <path> is a Kubernetes API path, e.g. /api/v1/pods or /apis/apps/v1/deployments.',
        '  --search/--page/--page-size enable server-side paging & keyword match (proxy search mode).',
        '  --filter/--sort/--limit are applied client-side on list items.',
        '  DELETE asks for confirmation unless --yes is given.',
      ].join('\n'),
    );
  c.allowUnknownOption();
  c.action(async (cluster: string, path: string) => {
    const opts = mergedOpts(c) as K8sOpts;
    const auth = resolveAuth(opts);
    const method = (opts.method ?? 'GET').toUpperCase();
    const useSearch = Boolean(opts.search || opts.page || opts.pageSize);
    const query: Record<string, string | number | boolean | undefined> = {
      namespace: opts.namespace,
      keywords: opts.search,
      search: useSearch ? 'true' : undefined,
      pageNum: opts.page,
      pageSize: opts.pageSize,
    };
    let body: unknown;
    if (opts.data) {
      try {
        body = JSON.parse(opts.data);
      } catch {
        throw E.badRequest('--data is not valid JSON');
      }
    }
    if (method === 'DELETE') await confirmAction(`DELETE ${path} on cluster "${cluster}"?`, opts);
    const proxyPath = `/proxy/${encodeURIComponent(cluster)}/k8s${path.startsWith('/') ? path : `/${path}`}`;
    const res = await request(auth, proxyPath, { method, query, body, raw: true });

    if (opts.format === 'json') {
      printData(res, { opts });
      return;
    }
    if (res && typeof res === 'object' && Array.isArray((res as { items?: unknown[] }).items)) {
      let items = (res as { items: unknown[] }).items;
      const groups = parseFilters(opts.filter);
      items = items.filter((i) => matchGroups(i, groups));
      if (opts.search && !useSearch) items = items.filter((i) => matchSearch(i, opts.search!));
      items = applySort(items, opts.sort);
      const limit = numOpt(opts.limit, '--limit');
      if (limit !== undefined) items = items.slice(0, limit);
      printData(items, { opts, defaultColumns: ['metadata.name', 'metadata.namespace', 'metadata.creationTimestamp'] });
      return;
    }
    printData(res, { opts });
  });
}
