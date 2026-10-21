import { Command } from 'commander';
import * as readline from 'node:readline';
import { decodeTokenProfile, ensureToken, loginForToken, resolveAuth } from './client.js';
import { E } from './errors.js';
import { mergedOpts } from './module.js';
import { printData } from './output.js';
import { clearConfig, loadConfig, saveConfig } from './store.js';

function question(query: string, mute: boolean): Promise<string> {
  return new Promise((resolve, reject) => {
    const rl = readline.createInterface({ input: process.stdin, output: process.stdout, terminal: true });
    if (mute) {
      const rlAny = rl as unknown as { _writeToOutput: (s: string) => void; output: NodeJS.WritableStream };
      rlAny._writeToOutput = (s: string) => {
        if (s.includes(query)) rlAny.output.write(s);
        else if (s === '\n' || s === '\r\n') rlAny.output.write(s);
        else rlAny.output.write('*');
      };
    }
    rl.question(query, (answer) => {
      rl.close();
      if (mute) process.stdout.write('\n');
      resolve(answer.trim());
    });
    rl.on('error', reject);
  });
}

export function registerAuth(program: Command): void {
  program
    .command('login')
    .description('log in to a KubePi server and cache credentials')
    .option('-s, --server <url>', 'KubePi server URL, e.g. http://127.0.0.1:2019')
    .option('-u, --username <name>', 'username or email')
    .option('-p, --password <password>', 'password')
    .addHelpText('after', '\nCredentials are cached in ~/.config/kubepi/config.json (mode 0600).\nEnv vars: KUBEPI_URL, KUBEPI_USERNAME, KUBEPI_PASSWORD, KUBEPI_TOKEN.')
    .action(async (_options: unknown, cmd: Command) => {
      const opts = mergedOpts(cmd) as { server?: string; username?: string; password?: string };
      const cfg = loadConfig();
      let server = opts.server || process.env.KUBEPI_URL || cfg.server;
      let username = opts.username || process.env.KUBEPI_USERNAME || cfg.username;
      let password = opts.password || process.env.KUBEPI_PASSWORD || cfg.password;
      const tty = process.stdin.isTTY;
      if (!server && tty) server = await question('KubePi server URL: ', false);
      if (!username && tty) username = await question('Username: ', false);
      if (!password && tty) password = await question('Password: ', true);
      if (!server) throw E.missing('--server');
      if (!username) throw E.missing('--username');
      if (!password) throw E.missing('--password');
      const auth = { server: server.replace(/\/+$/, ''), username, password };
      const token = await loginForToken(auth);
      saveConfig({ server: auth.server, username, password, token, tokenSavedAt: Date.now() });
      const profile = decodeTokenProfile(token);
      console.log(`logged in as ${(profile?.name as string) ?? username} (${auth.server})`);
    });

  program
    .command('logout')
    .description('clear cached credentials')
    .action(() => {
      clearConfig();
      console.log('logged out');
    });

  program
    .command('profile')
    .description('show the current login profile')
    .action(async (_options: unknown, cmd: Command) => {
      const opts = mergedOpts(cmd);
      const auth = resolveAuth(opts);
      const token = await ensureToken(auth);
      const profile = decodeTokenProfile(token);
      if (!profile) throw E.authExpired();
      const { exp, iat, nbf, ...rest } = profile as Record<string, unknown>;
      printData(rest, { opts });
    });
}
