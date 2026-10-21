import { Command } from 'commander';
import { cpSync, existsSync, statSync } from 'node:fs';
import { homedir } from 'node:os';
import { isAbsolute, join, relative, resolve, sep } from 'node:path';
import { fileURLToPath } from 'node:url';
import { E } from './errors.js';

const SKILL_NAME = 'kubepi-cli';
const AGENT_TARGETS = {
  'claude-code': { label: 'Claude Code', dir: '.claude/skills' },
  cursor: { label: 'Cursor', dir: '.cursor/skills' },
  'cherry-studio': { label: 'Cherry Studio', dir: '.cherrystudio/skills' },
  codex: { label: 'Codex', dir: '.agents/skills' },
  opencode: { label: 'OpenCode', dir: '.config/opencode/skills' },
  vscode: { label: 'VS Code', dir: '.copilot/skills' },
  antigravity: { label: 'Antigravity', dir: '.gemini/antigravity/skills' },
  gemini: { label: 'Gemini', dir: '.gemini/skills' },
  'kimi-code': { label: 'Kimi Code', dir: '.kimi-code/skills' },
  zcode: { label: 'ZCode', dir: '.zcode/skills' },
} as const;

type AgentName = keyof typeof AGENT_TARGETS;

interface AddSkillOptions {
  output?: string;
}

interface SkillTarget {
  label: string;
  destDir: string;
}

const AGENT_NAMES = Object.keys(AGENT_TARGETS) as AgentName[];

function resolveSkillSource(): string {
  const sourceDir = fileURLToPath(new URL(`../skills/${SKILL_NAME}/`, import.meta.url));
  if (!existsSync(join(sourceDir, 'SKILL.md'))) {
    throw new Error(`bundled skill not found: ${sourceDir}; reinstall the complete kubepi-cli package`);
  }
  return sourceDir;
}

function resolveAgents(agent: string): AgentName[] {
  const normalized = agent.toLowerCase();
  if (normalized === 'all') return AGENT_NAMES;
  if (normalized in AGENT_TARGETS) return [normalized as AgentName];
  throw E.badRequest(`unsupported agent: ${agent}; available: ${AGENT_NAMES.join(', ')}, all`);
}

function resolveAgentTarget(agent: AgentName, homeDir: string): SkillTarget {
  const target = AGENT_TARGETS[agent];
  const destDir = agent === 'kimi-code' && process.env.KIMI_CODE_HOME
    ? resolve(join(process.env.KIMI_CODE_HOME, 'skills', SKILL_NAME))
    : resolve(join(homeDir, target.dir, SKILL_NAME));
  return { label: target.label, destDir };
}

function resolveOutputDir(output: string, homeDir: string): string {
  const trimmed = output.trim();
  if (!trimmed) throw E.badRequest('output path cannot be empty');
  const expanded = trimmed === '~'
    ? homeDir
    : /^~[\\/]/.test(trimmed)
      ? join(homeDir, trimmed.slice(2))
      : trimmed;
  const outputDir = resolve(expanded);
  if (existsSync(outputDir) && !statSync(outputDir).isDirectory()) {
    throw E.badRequest(`output path must be a directory: ${outputDir}`);
  }
  return outputDir;
}

function assertSkillDestination(sourceDir: string, destDir: string): void {
  const relativeDestination = relative(resolve(sourceDir), resolve(destDir));
  const parentPrefix = '..' + sep;
  const destinationIsWithinSource = relativeDestination === ''
    || (!isAbsolute(relativeDestination)
      && relativeDestination !== '..'
      && !relativeDestination.startsWith(parentPrefix));
  if (destinationIsWithinSource) {
    throw E.badRequest(`skill destination cannot be the source directory or its child: ${destDir}`);
  }
  let existingAncestor = resolve(destDir);
  while (!existsSync(existingAncestor)) {
    const parent = resolve(existingAncestor, '..');
    if (parent === existingAncestor) break;
    existingAncestor = parent;
  }
  if (!statSync(existingAncestor).isDirectory()) {
    throw E.badRequest(`skill destination must be a directory: ${existingAncestor}`);
  }
}

export function registerAddSkill(program: Command, homeDir: string = homedir()): void {
  program
    .command('add-skill')
    .description('install the bundled Agent skill or export it to a custom directory')
    .argument('[agent]', `target Agent (${AGENT_NAMES.join('|')}|all)`)
    .option('-o, --output <path>', 'export the bundled skill to a custom directory')
    .action((agent?: string, options: AddSkillOptions = {}) => {
      if (options.output !== undefined && agent) {
        throw E.badRequest('agent and --output cannot be used together');
      }
      if (options.output === undefined && !agent) {
        throw E.missing(`agent (${AGENT_NAMES.join('|')}|all) or --output`);
      }

      const sourceDir = resolveSkillSource();
      const targets: SkillTarget[] = options.output !== undefined
        ? [{ label: 'exported', destDir: join(resolveOutputDir(options.output, homeDir), SKILL_NAME) }]
        : resolveAgents(agent!).map((name) => resolveAgentTarget(name, homeDir));

      for (const target of targets) assertSkillDestination(sourceDir, target.destDir);
      for (const target of targets) {
        cpSync(sourceDir, target.destDir, { recursive: true, dereference: true, force: true });
        if (target.label === 'exported') console.log(`exported ${SKILL_NAME} skill to: ${target.destDir}`);
        else console.log(`installed ${SKILL_NAME} skill for ${target.label}: ${target.destDir}`);
      }
    });
}
