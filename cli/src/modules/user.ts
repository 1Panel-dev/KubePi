import { api } from '../client.js';
import { E } from '../errors.js';
import { deepMerge, registerModule, type Opts } from '../module.js';
import { collect, str } from '../shared.js';

export function registerUser(program: Parameters<typeof registerModule>[0]): void {
  registerModule(program, {
    name: 'user',
    description: 'manage users (用户)',
    basePath: '/users',
    searchPath: '/users/search',
    defaultColumns: ['name', 'nickName', 'email', 'type', 'isAdmin', 'createAt'],
    createOptions: (c) => {
      c.option('--name <name>', 'username');
      c.option('--nick-name <name>', 'nick name');
      c.option('--email <email>', 'email');
      c.option('--password <password>', 'login password (required)');
      c.option('--language <lang>', 'zh-CN or en-US');
      c.option('--is-admin', 'grant administrator');
      c.option('--role <role>', 'role name, repeatable', collect);
    },
    updateOptions: (c) => {
      c.option('--nick-name <name>', 'nick name');
      c.option('--email <email>', 'email');
      c.option('--language <lang>', 'zh-CN or en-US');
      c.option('--is-admin <bool>', 'true or false');
      c.option('--role <role>', 'replace roles, repeatable', collect);
      c.option('--password <password>', 'new password (requires --old-password)');
      c.option('--old-password <password>', 'current password');
    },
    buildCreate: (fields, opts) => {
      const name = str(opts.name) ?? str(fields.name);
      if (!name) throw E.missing('--name');
      const password = str(opts.password) ?? str(fields.password) ?? str((fields.authenticate as Record<string, unknown> | undefined)?.password);
      if (!password) throw E.missing('--password');
      const roles = (opts.role as string[] | undefined) ?? (fields.roles as string[] | undefined) ?? [];
      return deepMerge(fields, {
        apiVersion: 'v1',
        kind: 'User',
        name,
        nickName: str(opts.nickName) ?? str(fields.nickName) ?? name,
        email: str(opts.email) ?? str(fields.email) ?? '',
        language: str(opts.language) ?? str(fields.language),
        isAdmin: opts.isAdmin === true || fields.isAdmin === true,
        type: 'LOCAL',
        authenticate: { password },
        roles,
      });
    },
    buildUpdate: (fields, current, opts) => {
      if (opts.password) {
        if (!opts.oldPassword) throw E.missing('--old-password (required to change password)');
        return { password: opts.password, oldPassword: opts.oldPassword };
      }
      const patch: Record<string, unknown> = deepMerge(current ?? {}, fields) as Record<string, unknown>;
      if (opts.nickName) patch.nickName = opts.nickName;
      if (opts.email) patch.email = opts.email;
      if (opts.language) patch.language = opts.language;
      if (opts.isAdmin !== undefined) patch.isAdmin = opts.isAdmin === true || opts.isAdmin === 'true';
      const roles = (opts.role as string[] | undefined) ?? (fields.roles as string[] | undefined);
      if (roles) patch.roles = roles;
      delete patch.password;
      delete patch.oldPassword;
      return patch;
    },
    actions: [
      {
        name: 'reset-password <name>',
        description: 'admin resets a user password (requires KubePi >= v2.0)',
        options: (c) => c.option('--password <password>', 'new password'),
        run: async (args, opts, auth) => {
          if (!opts.password) throw E.missing('--password');
          return api.put(auth, `/users/${encodeURIComponent(String(args[0]))}/password`, { password: opts.password });
        },
      },
    ],
    help: [
      'Examples:',
      '  kubepi user create --name=tom --email=tom@example.com --password=Secret123 --role=dev',
      '  kubepi user update tom --role=ops --role=dev',
      '  kubepi user update tom --password=NewSecret --old-password=Secret123',
    ].join('\n'),
  });
}
