import { Command } from 'commander';
import { registerAuth } from './auth.js';
import { hoistGlobalOptions } from './argv.js';
import { CliError, errorMessage } from './errors.js';
import { registerChart } from './modules/chart.js';
import { registerCluster } from './modules/cluster.js';
import { registerImageRepo } from './modules/imagerepo.js';
import { registerK8s } from './modules/k8s.js';
import { registerLdap } from './modules/ldap.js';
import { registerMonitor } from './modules/monitor.js';
import { registerRole } from './modules/role.js';
import { registerSystem } from './modules/system.js';
import { registerUser } from './modules/user.js';
import { registerAddSkill } from './skill.js';
import { collect } from './shared.js';
import pkg from '../package.json';

const MODULES_HELP = `
Modules:
  login / logout / profile   authentication
  cluster                    clusters (members/clusterroles/repos/namespaces via actions)
  clustermember              cluster members        [--cluster]
  clusterrole                K8s ClusterRoles       [--cluster]
  clusterrepo                cluster chart repos    [--cluster]
  user                       users
  role                       roles
  imagerepo                  image registries
  chartrepo                  helm chart repos       [--cluster]
  chart                      browse/install charts  [--cluster]
  ldap                       LDAP settings
  sso                        SSO settings
  monitor                    grafana & metrics      (KubePi >= v2.0)
  log                        login / operation logs
  k8s                        generic Kubernetes API passthrough
  add-skill                  install the bundled Agent skill for supported AI agents

Shorthand:
  kubepi <module>                     list
  kubepi <module> <name>              detail
  kubepi <module> create --field=v    create (or --data='JSON')
  kubepi <module> update <name> --f=v update
  kubepi <module> delete <name>       delete (asks unless --yes)
  kubepi <module> help                module help

Examples:
  kubepi login -s http://127.0.0.1:2019 -u admin -p kubepi
  kubepi user --search=tom --format=json
  kubepi cluster create --name=dev --mode=configFile --kubeconfig=~/.kube/config
  kubepi k8s dev /api/v1/pods --namespace=default
  kubepi add-skill kimi-code
`;

function buildProgram(): Command {
  const program = new Command();
  program
    .name('kubepi')
    .description('KubePi CLI - operate KubePi from the terminal; friendly to AI agents')
    .version(pkg.version)
    .option('--server <url>', 'KubePi server URL (env: KUBEPI_URL)')
    .option('--username <name>', 'username (env: KUBEPI_USERNAME)')
    .option('--password <password>', 'password (env: KUBEPI_PASSWORD)')
    .option('--token <token>', 'JWT token (env: KUBEPI_TOKEN)')
    .option('--format <format>', 'output format: table (default) | json')
    .option('--pick <fields>', 'comma-separated fields to pick, dot paths supported')
    .option('--filter <expr>', 'filter expr, repeatable (OR); comma-separated inside means AND', collect, [])
    .option('--search <keyword>', 'fuzzy search keyword')
    .option('--sort <keys>', 'sort keys, e.g. name_asc,createAt_desc')
    .option('--page <n>', 'page number')
    .option('--page-size <n>', 'page size')
    .option('--all', 'fetch all pages')
    .option('--limit <n>', 'limit result count')
    .option('--yes', 'skip confirmation prompts')
    .option('--data <json>', 'request body as JSON')
    .showSuggestionAfterError()
    .addHelpText('after', MODULES_HELP);
  
  registerAuth(program);
  registerAddSkill(program);
  registerCluster(program);
  registerUser(program);
  registerRole(program);
  registerImageRepo(program);
  registerChart(program);
  registerLdap(program);
  registerSystem(program);
  registerMonitor(program);
  registerK8s(program);
  return program;
}

async function main(): Promise<void> {
  const program = buildProgram();
  try {
    await program.parseAsync(hoistGlobalOptions(process.argv));
  } catch (err) {
    if (err instanceof CliError) {
      console.error(`${err.code}: ${err.message}`);
      if (err.hint) console.error(`hint: ${err.hint}`);
      process.exit(1);
    }
    console.error(`E5001: ${errorMessage(err)}`);
    process.exit(1);
  }
}

void main();
