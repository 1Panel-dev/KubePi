import { Command } from 'commander';
import { api, resolveAuth, type GlobalOpts, type ResolvedAuth } from './client.js';
import { CliError, E } from './errors.js';
import { printData, setPath } from './output.js';
import { applyClientQuery, buildConditions, parseFilters } from './query.js';
import { confirmAction } from './shared.js';

export { confirmAction };

export interface ActionDef {
  /** Usage string, e.g. 'sync <name>'. */
  name: string;
  description: string;
  options?: (cmd: Command) => void;
  /** Extra help lines appended to the action help. */
  help?: string;
  run: (
    args: unknown[],
    opts: GlobalOpts & Record<string, unknown>,
    auth: ResolvedAuth,
    rawArgs: string[],
  ) => Promise<unknown>;
  printColumns?: string[];
}

export interface ModuleDef {
  name: string;
  description: string;
  basePath?: string | ((opts: Record<string, unknown>) => string);
  searchPath?: string | ((opts: Record<string, unknown>) => string);
  /** Option keys that must be present (e.g. ['cluster']). */
  requires?: string[];
  defaultColumns?: string[];
  noList?: boolean;
  noGet?: boolean;
  noCreate?: boolean;
  noUpdate?: boolean;
  noDelete?: boolean;
  /** Skip the GET-before-PUT prefill for update. */
  noPrefill?: boolean;
  createOptions?: (cmd: Command) => void;
  updateOptions?: (cmd: Command) => void;
  buildCreate?: (fields: Record<string, unknown>, opts: GlobalOpts & Record<string, unknown>) => unknown | Promise<unknown>;
  buildUpdate?: (
    fields: Record<string, unknown>,
    current: unknown,
    opts: GlobalOpts & Record<string, unknown>,
  ) => unknown;
  actions?: ActionDef[];
  /** Extra help lines appended to the module help. */
  help?: string;
  /** Module-level options available on list/get and inherited by actions. */
  moduleOptions?: (cmd: Command) => void;
  /** Extra query params for the plain-GET list path. */
  listQuery?: (opts: Opts) => Record<string, string | number | boolean | undefined>;
  /** Post-process output data (e.g. strip sensitive fields). */
  sanitize?: (data: unknown) => unknown;
}

export type Opts = GlobalOpts & Record<string, unknown>;

export function resolvePath(p: string | ((opts: Opts) => string) | undefined, opts: Opts): string {
  if (p === undefined) throw E.unknownModule('(no path configured)');
  return typeof p === 'function' ? p(opts) : p;
}

function requireOpts(def: ModuleDef, opts: Opts): void {
  for (const key of def.requires ?? []) {
    if (opts[key] === undefined || opts[key] === '') throw E.missing(`--${key.replace(/[A-Z]/g, (c) => '-' + c.toLowerCase())}`);
  }
}

function coerce(raw: string): unknown {
  const trimmed = raw.trim();
  if (trimmed === '') return '';
  try {
    return JSON.parse(trimmed);
  } catch {
    return raw;
  }
}

/** Merge --data JSON and pass-through --field=value pairs (dot paths supported). */
export function collectFields(opts: Opts, args: string[]): Record<string, unknown> {
  const fields: Record<string, unknown> = {};
  if (opts.data) {
    let parsed: unknown;
    try {
      parsed = JSON.parse(opts.data);
    } catch {
      throw E.badRequest('--data is not valid JSON');
    }
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) Object.assign(fields, parsed);
    else throw E.badRequest('--data must be a JSON object');
  }
  for (let i = 0; i < args.length; i++) {
    const a = args[i];
    if (!a.startsWith('--')) continue;
    const eq = a.indexOf('=');
    if (eq > 2) {
      setPath(fields, a.slice(2, eq), coerce(a.slice(eq + 1)));
    } else if (i + 1 < args.length && !args[i + 1].startsWith('--')) {
      setPath(fields, a.slice(2), coerce(args[++i]));
    } else {
      setPath(fields, a.slice(2), true);
    }
  }
  return fields;
}

function isPlainObject(v: unknown): v is Record<string, unknown> {
  return typeof v === 'object' && v !== null && !Array.isArray(v);
}

export function deepMerge(base: unknown, patch: unknown): unknown {
  if (!isPlainObject(base) || !isPlainObject(patch)) return patch === undefined ? base : patch;
  const out: Record<string, unknown> = { ...base };
  for (const [k, v] of Object.entries(patch)) {
    out[k] = k in out ? deepMerge(out[k], v) : v;
  }
  return out;
}

interface PageLike {
  items?: unknown[];
  total?: number;
}

function unwrapItems(res: unknown): { items: unknown[]; total?: number } {
  if (Array.isArray(res)) return { items: res };
  if (isPlainObject(res) && Array.isArray((res as PageLike).items)) {
    return { items: (res as PageLike).items!, total: (res as PageLike).total };
  }
  return { items: [] };
}

export async function listItems(def: ModuleDef, opts: Opts, auth: ResolvedAuth): Promise<unknown[]> {
  const groups = parseFilters(opts.filter);
  const searchPath = def.searchPath ? resolvePath(def.searchPath, opts) : undefined;
  const conditions = searchPath ? buildConditions(groups, opts.search) : null;

  let items: unknown[];
  if (searchPath && conditions) {
    const body = { conditions };
    if (opts.all) {
      items = [];
      const size = 100;
      for (let page = 1; ; page++) {
        const { items: pageItems, total } = unwrapItems(await api.post(auth, searchPath, body, { pageNum: page, pageSize: size }));
        items.push(...pageItems);
        if (pageItems.length < size || (total !== undefined && items.length >= total)) break;
      }
    } else {
      const query =
        opts.page || opts.pageSize
          ? { pageNum: opts.page ?? '1', pageSize: opts.pageSize ?? '10' }
          : undefined;
      items = unwrapItems(await api.post(auth, searchPath, body, query)).items;
    }
  } else {
    if (!def.basePath) throw E.unknownModule(def.name);
    items = unwrapItems(await api.get(auth, resolvePath(def.basePath, opts), def.listQuery?.(opts))).items;
  }

  return applyClientQuery(items, opts);
}

/**
 * Merge options along the command chain, program first, subcommand last;
 * at each level only defined values win. This lets module options be given
 * either before or after an action subcommand, and lets global options
 * coexist with same-named local declarations.
 */
export function mergedOpts(cmd: Command): Opts {
  const chain: Command[] = [];
  for (let c: Command | null = cmd; c; c = c.parent as Command | null) chain.unshift(c);
  const out: Record<string, unknown> = {};
  for (const c of chain) {
    for (const [k, v] of Object.entries(c.opts() as Record<string, unknown>)) {
      if (v !== undefined) out[k] = v;
    }
  }
  return out as Opts;
}

export function registerModule(program: Command, def: ModuleDef): void {
  const sanitized = (data: unknown): unknown => {
    if (!def.sanitize) return data;
    return Array.isArray(data) ? data.map((i) => def.sanitize!(i)) : def.sanitize(data);
  };
  const mod = program.command(def.name).description(def.description);
  if (def.help) mod.addHelpText('after', `\n${def.help}`);
  def.moduleOptions?.(mod);
  mod.command('help').description('show help for this module').action(() => mod.outputHelp());

  if (!def.noList || !def.noGet) {
    mod.argument('[name]').action(async (name?: string) => {
      const opts = mergedOpts(mod);
      const auth = resolveAuth(opts);
      requireOpts(def, opts);
      if (name !== undefined) {
        if (def.noGet || !def.basePath) throw E.unknownModule(name);
        printData(sanitized(await api.get(auth, `${resolvePath(def.basePath, opts)}/${encodeURIComponent(name)}`)), { opts });
      } else {
        if (def.noList) throw E.missing('name');
        printData(sanitized(await listItems(def, opts, auth)), { opts, defaultColumns: def.defaultColumns });
      }
    });
  }

  if (!def.noCreate) {
    const c = mod.command('create').description(`create a ${def.name}`);
    c.allowUnknownOption();
    c.argument('[extras...]', 'extra --field=value pairs');
    def.createOptions?.(c);
    c.action(async () => {
      if (c.args.includes('help')) {
        c.outputHelp();
        return;
      }
      const opts = mergedOpts(c);
      const auth = resolveAuth(opts);
      requireOpts(def, opts);
      const fields = collectFields(opts, c.args);
      const body = def.buildCreate ? await def.buildCreate(fields, opts) : fields;
      printData(sanitized(await api.post(auth, resolvePath(def.basePath!, opts), body)), { opts });
    });
  }

  if (!def.noUpdate) {
    const c = mod.command('update <name>').description(`update a ${def.name}`);
    c.allowUnknownOption();
    def.updateOptions?.(c);
    c.action(async (name: string) => {
      if (name === 'help' || c.args.includes('help')) {
        c.outputHelp();
        return;
      }
      const opts = mergedOpts(c);
      const auth = resolveAuth(opts);
      requireOpts(def, opts);
      const fields = collectFields(opts, c.args);
      const base = resolvePath(def.basePath!, opts);
      const current = def.noPrefill ? undefined : await api.get(auth, `${base}/${encodeURIComponent(name)}`);
      const body = def.buildUpdate ? def.buildUpdate(fields, current, opts) : deepMerge(current ?? {}, fields);
      printData(sanitized(await api.put(auth, `${base}/${encodeURIComponent(name)}`, body)), { opts });
    });
  }

  if (!def.noDelete) {
    const c = mod.command('delete <name>').description(`delete a ${def.name}`);
    c.action(async (name: string) => {
      if (name === 'help') {
        c.outputHelp();
        return;
      }
      const opts = mergedOpts(c);
      const auth = resolveAuth(opts);
      requireOpts(def, opts);
      await confirmAction(`Delete ${def.name} "${name}"?`, opts);
      await api.del(auth, `${resolvePath(def.basePath!, opts)}/${encodeURIComponent(name)}`);
      console.log(`deleted ${def.name} "${name}"`);
    });
  }

  for (const action of def.actions ?? []) {
    const c = mod.command(action.name).description(action.description);
    c.allowUnknownOption();
    c.argument('[extras...]', 'extra --field=value pairs');
    action.options?.(c);
    if (action.help) c.addHelpText('after', `\n${action.help}`);
    c.action(async (...args: unknown[]) => {
      const cmdObj = args[args.length - 1] as Command;
      if (cmdObj.args.includes('help') || cmdObj.processedArgs.includes('help')) {
        cmdObj.outputHelp();
        return;
      }
      const opts = mergedOpts(cmdObj);
      const auth = resolveAuth(opts);
      requireOpts(def, opts);
      const result = await action.run(cmdObj.processedArgs, opts, auth, cmdObj.args);
      const out = Array.isArray(result) ? applyClientQuery(result, opts) : result;
      if (out !== undefined) printData(out, { opts, defaultColumns: action.printColumns });
    });
  }
}
