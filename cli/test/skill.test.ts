import { afterEach, beforeEach, describe, expect, test } from 'bun:test';
import { Command } from 'commander';
import { existsSync, mkdirSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { registerAddSkill } from '../src/skill.js';

const sourceDir = resolve('skills/kubepi-cli');
const sourceSkill = readFileSync(join(sourceDir, 'SKILL.md'), 'utf8');
const legacyAgentTargets = [
  ['claude-code', '.claude/skills/kubepi-cli'],
  ['cursor', '.cursor/skills/kubepi-cli'],
  ['cherry-studio', '.cherrystudio/skills/kubepi-cli'],
  ['codex', '.agents/skills/kubepi-cli'],
  ['opencode', '.config/opencode/skills/kubepi-cli'],
  ['vscode', '.copilot/skills/kubepi-cli'],
  ['antigravity', '.gemini/antigravity/skills/kubepi-cli'],
  ['gemini', '.gemini/skills/kubepi-cli'],
] as const;
const kimiTarget = '.kimi-code/skills/kubepi-cli';
const zcodeTarget = '.zcode/skills/kubepi-cli';

let dir: string;
let home: string;
let originalKimiCodeHome: string | undefined;

beforeEach(() => {
  dir = mkdtempSync(join(tmpdir(), 'kubepi-cli-skill-test-'));
  home = join(dir, 'home');
  mkdirSync(home, { recursive: true });
  originalKimiCodeHome = process.env.KIMI_CODE_HOME;
  delete process.env.KIMI_CODE_HOME;
});

afterEach(() => {
  if (originalKimiCodeHome === undefined) delete process.env.KIMI_CODE_HOME;
  else process.env.KIMI_CODE_HOME = originalKimiCodeHome;
  rmSync(dir, { recursive: true, force: true });
});

async function run(args: string[]): Promise<void> {
  const program = new Command();
  registerAddSkill(program, home);
  await program.parseAsync(args, { from: 'user' });
}

function expectCopied(destDir: string): void {
  expect(readFileSync(join(destDir, 'SKILL.md'), 'utf8')).toBe(sourceSkill);
}

describe('add-skill', () => {
  for (const [agent, target] of legacyAgentTargets) {
    test(`installs the bundled skill for ${agent}`, async () => {
      await run(['add-skill', agent]);
      expectCopied(join(home, target));
    });
  }

  test('installs for Kimi Code in the default data root', async () => {
    await run(['add-skill', 'kimi-code']);
    expectCopied(join(home, kimiTarget));
  });

  test('installs for Kimi Code in KIMI_CODE_HOME', async () => {
    const kimiHome = join(dir, 'custom-kimi');
    process.env.KIMI_CODE_HOME = kimiHome;
    await run(['add-skill', 'kimi-code']);
    expectCopied(join(kimiHome, 'skills/kubepi-cli'));
    expect(existsSync(join(home, kimiTarget))).toBe(false);
  });

  test('treats an empty KIMI_CODE_HOME as unset', async () => {
    process.env.KIMI_CODE_HOME = '';
    await run(['add-skill', 'kimi-code']);
    expectCopied(join(home, kimiTarget));
  });

  test('installs for ZCode independently of KIMI_CODE_HOME and agent case', async () => {
    process.env.KIMI_CODE_HOME = join(dir, 'custom-kimi');
    await run(['add-skill', 'ZCODE']);
    expectCopied(join(home, zcodeTarget));
    expect(existsSync(join(dir, 'custom-kimi'))).toBe(false);
  });

  test('installs all ten agent targets', async () => {
    await run(['add-skill', 'all']);
    for (const [, target] of legacyAgentTargets) expectCopied(join(home, target));
    expectCopied(join(home, kimiTarget));
    expectCopied(join(home, zcodeTarget));
  });

  test('exports to a custom directory with spaces and non-ASCII characters', async () => {
    const output = join(dir, '导出 skills');
    await run(['add-skill', '--output', output]);
    expectCopied(join(output, 'kubepi-cli'));
  });

  test('expands tilde in a custom output directory', async () => {
    await run(['add-skill', '--output', '~/exported']);
    expectCopied(join(home, 'exported/kubepi-cli'));
  });

  test('overwrites the skill on reinstall but preserves extra user files', async () => {
    const destDir = join(home, zcodeTarget);
    await run(['add-skill', 'zcode']);
    writeFileSync(join(destDir, 'SKILL.md'), 'outdated');
    writeFileSync(join(destDir, 'user-note.md'), 'keep');

    await run(['add-skill', 'zcode']);

    expectCopied(destDir);
    expect(readFileSync(join(destDir, 'user-note.md'), 'utf8')).toBe('keep');
  });

  test('rejects invalid or conflicting arguments without copying', async () => {
    await expect(run(['add-skill', 'not-an-agent'])).rejects.toMatchObject({ code: 'E2004' });
    await expect(run(['add-skill'])).rejects.toMatchObject({ code: 'E2003' });
    await expect(run(['add-skill', 'kimi-code', '--output', join(dir, 'out')])).rejects.toMatchObject({ code: 'E2004' });
    await expect(run(['add-skill', '--output', '   '])).rejects.toMatchObject({ code: 'E2004' });
    expect(existsSync(join(home, '.kimi-code'))).toBe(false);
  });

  test('rejects an output root that is a regular file', async () => {
    const output = join(dir, 'not-a-directory');
    writeFileSync(output, 'unchanged');
    await expect(run(['add-skill', '--output', output])).rejects.toMatchObject({ code: 'E2004' });
    expect(readFileSync(output, 'utf8')).toBe('unchanged');
  });

  test('rejects exporting into the skill source or its child directory', async () => {
    await expect(run(['add-skill', '--output', resolve('skills')])).rejects.toMatchObject({ code: 'E2004' });
    await expect(run(['add-skill', '--output', sourceDir])).rejects.toMatchObject({ code: 'E2004' });
  });

  test('validates every destination before installing all targets', async () => {
    const zcodeParent = join(home, '.zcode/skills');
    mkdirSync(zcodeParent, { recursive: true });
    writeFileSync(join(zcodeParent, 'kubepi-cli'), 'conflict');

    await expect(run(['add-skill', 'all'])).rejects.toMatchObject({ code: 'E2004' });

    for (const [, target] of legacyAgentTargets) expect(existsSync(join(home, target))).toBe(false);
    expect(existsSync(join(home, kimiTarget))).toBe(false);
    expect(readFileSync(join(zcodeParent, 'kubepi-cli'), 'utf8')).toBe('conflict');
  });

  test('validates existing ancestors before installing all targets', async () => {
    const zcodeHome = join(home, '.zcode');
    writeFileSync(zcodeHome, 'conflict');

    await expect(run(['add-skill', 'all'])).rejects.toMatchObject({ code: 'E2004' });

    for (const [, target] of legacyAgentTargets) expect(existsSync(join(home, target))).toBe(false);
    expect(existsSync(join(home, kimiTarget))).toBe(false);
    expect(readFileSync(zcodeHome, 'utf8')).toBe('conflict');
  });
});
