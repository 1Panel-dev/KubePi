import { JSONPath } from 'jsonpath-plus';
const quotedKey = /['"]/;
const quotedMatch = /[^."']+|"([^"]*)"|'([^']*)'/g;
import Vue from 'vue';


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

export function set(obj, path, value) {
  let ptr = obj;
  let parts;

  if (!ptr) {
    return;
  }

  if ( path.match(quotedKey) ) {
    // Path with quoted section
    parts = path.match(quotedMatch).map(x => x.replace(/['"]/g, ''));
  } else {
    // Regular path
    parts = path.split('.');
  }

  for (let i = 0; i < parts.length; i++) {
    const key = parts[i];

    if ( i === parts.length - 1 ) {
      Vue.set(ptr, key, value);
    } else if ( !ptr[key] ) {
      // Make sure parent keys exist
      Vue.set(ptr, key, {});
    }

    ptr = ptr[key];
  }

  return obj;
}
