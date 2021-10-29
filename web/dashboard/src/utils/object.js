import { JSONPath } from 'jsonpath-plus';
const quotedKey = /['"]/;


export function get(obj, path) {
  if ( path.startsWith('$') ) {
    try {
      return JSONPath({
        path,
        json:        obj,
        preventEval: true,
        wrap:        false,
      });
    } catch (e) {
      console.log('JSON Path error', e, path, obj); // eslint-disable-line no-console

      return '(JSON Path err)';
    }
  }

  let parts;

  if ( path.match(quotedKey) ) {
    // Path with quoted section
    parts = path.match(/[^."']+|"([^"]*)"|'([^']*)'/g).map(x => x.replace(/['"]/g, ''));
  } else {
    // Regular path
    parts = path.split('.');
  }

  for (let i = 0; i < parts.length; i++) {
    if (!obj) {
      return;
    }

    obj = obj[parts[i]];
  }

  return obj;
}
