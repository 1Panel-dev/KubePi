export class CliError extends Error {
  readonly code: string;
  readonly hint?: string;

  constructor(code: string, message: string, hint?: string) {
    super(message);
    this.name = 'CliError';
    this.code = code;
    this.hint = hint;
  }
}

export const E = {
  notLoggedIn: () =>
    new CliError(
      'E1001',
      'not logged in: missing KubePi server credentials',
      'run `kubepi login -s <server> -u <username> -p <password>`, or set KUBEPI_URL / KUBEPI_USERNAME / KUBEPI_PASSWORD',
    ),
  loginFailed: (msg: string) => new CliError('E1004', `login failed: ${msg}`, 'check username/password, then run `kubepi login` again'),
  authExpired: () => new CliError('E1004', 'authentication expired or token invalid', 'run `kubepi login` again'),
  mfaRequired: () =>
    new CliError(
      'E1005',
      'this account has MFA enabled; JWT login is rejected by the server (KubePi >= v2.0)',
      'use an account without MFA, or disable MFA for the account in KubePi',
    ),
  unknownModule: (name: string) => new CliError('E2001', `unknown module or command: ${name}`, 'run `kubepi help` to list available modules'),
  notFound: (what: string) =>
    new CliError('E2002', `not found: ${what}`, 'check the name; note some endpoints (monitor, user reset-password) require KubePi >= v2.0'),
  missing: (what: string) => new CliError('E2003', `missing required argument or option: ${what}`, 'run `kubepi <module> help` or `kubepi <module> <action> help`'),
  forbidden: (msg: string) => new CliError('E2006', `permission denied: ${msg}`, 'check the RBAC permissions of the current account'),
  invalid: (msg: string) => new CliError('E2004', msg),
  badRequest: (msg: string) => new CliError('E2004', msg),
  network: (msg: string) => new CliError('E5001', `request failed: ${msg}`, 'check the network and KubePi server status'),
} as const;

export function errorMessage(err: unknown): string {
  if (err instanceof Error) return err.message;
  return String(err);
}
