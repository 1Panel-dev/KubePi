import { chmodSync, existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { dirname, join } from 'node:path';

export interface StoredConfig {
  server?: string;
  username?: string;
  password?: string;
  token?: string;
  tokenSavedAt?: number;
}

function configDir(): string {
  return process.env.KUBEPI_CONFIG_DIR || join(homedir(), '.config', 'kubepi');
}

export function configFile(): string {
  return join(configDir(), 'config.json');
}

export function loadConfig(): StoredConfig {
  const file = configFile();
  if (!existsSync(file)) return {};
  try {
    return JSON.parse(readFileSync(file, 'utf8')) as StoredConfig;
  } catch {
    return {};
  }
}

export function saveConfig(cfg: StoredConfig): void {
  const file = configFile();
  mkdirSync(dirname(file), { recursive: true, mode: 0o700 });
  writeFileSync(file, JSON.stringify(cfg, null, 2), { mode: 0o600 });
  try {
    chmodSync(file, 0o600);
  } catch {
    // best effort on platforms without chmod
  }
}

export function patchConfig(patch: Partial<StoredConfig>): StoredConfig {
  const next = { ...loadConfig(), ...patch };
  saveConfig(next);
  return next;
}

export function clearConfig(): void {
  const file = configFile();
  if (existsSync(file)) rmSync(file);
}
