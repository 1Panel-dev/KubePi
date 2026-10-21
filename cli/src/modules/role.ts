import { E } from '../errors.js';
import { deepMerge, registerModule } from '../module.js';
import { str } from '../shared.js';

export function registerRole(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'role',
    description: 'manage roles (角色)',
    basePath: '/roles',
    searchPath: '/roles/search',
    defaultColumns: ['name', 'description', 'createAt'],
    createOptions: (c) => {
      c.option('--name <name>', 'role name');
      c.option('--description <text>', 'description');
      c.option('--rules <json>', 'rules JSON array, e.g. [{"resource":["clusters"],"verbs":["get","list"]}]');
    },
    updateOptions: (c) => {
      c.option('--description <text>', 'description');
      c.option('--rules <json>', 'rules JSON array');
    },
    buildCreate: (fields, opts) => {
      const name = str(opts.name) ?? str(fields.name);
      if (!name) throw E.missing('--name');
      let rules = fields.rules;
      if (opts.rules) {
        try {
          rules = JSON.parse(String(opts.rules));
        } catch {
          throw E.badRequest('--rules is not valid JSON');
        }
      }
      if (!Array.isArray(rules) || rules.length === 0) throw E.missing('--rules');
      return deepMerge(fields, {
        apiVersion: 'v1',
        kind: 'Role',
        name,
        description: str(opts.description) ?? str(fields.description) ?? '',
        rules,
      });
    },
    buildUpdate: (fields, current, opts) => {
      const patch = deepMerge(current ?? {}, fields) as Record<string, unknown>;
      if (opts.description) patch.description = opts.description;
      if (opts.rules) {
        try {
          patch.rules = JSON.parse(String(opts.rules));
        } catch {
          throw E.badRequest('--rules is not valid JSON');
        }
      }
      return patch;
    },
    help: [
      'Examples:',
      `  kubepi role create --name=dev --rules='[{"resource":["clusters","proxy"],"resourceNames":["dev"],"verbs":["get","list"]}]'`,
      '  kubepi role --filter=name:dev',
    ].join('\n'),
  });
}
