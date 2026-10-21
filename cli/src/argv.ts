const VALUE_FLAGS = new Set([
  '--server',
  '--username',
  '--password',
  '--token',
  '--format',
  '--pick',
  '--filter',
  '--search',
  '--sort',
  '--page',
  '--page-size',
  '--limit',
  '--data',
]);
const BOOL_FLAGS = new Set(['--all', '--yes']);

/**
 * Global options may appear anywhere on the command line (before or after
 * module/action names). Hoist them to the front so commander parses them
 * as program options regardless of position.
 */
export function hoistGlobalOptions(argv: string[]): string[] {
  const [node, script, ...rest] = argv;
  const hoisted: string[] = [];
  const remainder: string[] = [];
  for (let i = 0; i < rest.length; i++) {
    const a = rest[i];
    const eq = a.indexOf('=');
    const flag = eq > 0 ? a.slice(0, eq) : a;
    if (VALUE_FLAGS.has(flag)) {
      if (eq > 0) {
        hoisted.push(a);
      } else if (i + 1 < rest.length && !rest[i + 1].startsWith('--')) {
        hoisted.push(a, rest[++i]);
      } else {
        // No value given: leave the flag bare so commander reports the
        // missing argument instead of silently eating the next flag.
        hoisted.push(a);
      }
    } else if (BOOL_FLAGS.has(flag)) {
      hoisted.push(a);
    } else {
      remainder.push(a);
    }
  }
  return [node, script, ...hoisted, ...remainder];
}
