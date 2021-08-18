import {dateFormat, datetimeFormat} from "fit2cloud-ui/src/filters/time";

export function errorFormat(value) {
  if (value !== null) {
      let errItem = value;
      errItem = errItem.replace(/\\n/gi,'\n');
      errItem = errItem.replace(/\\u/gi,'%u');
      errItem = errItem.replace(/\\/gi,'');
      errItem = unescape(errItem)
      return errItem
  }
  return value
}

export function emailFormat(value) {
  let result = '';
  if (value.indexOf('@') === -1 || value.indexOf('.') === -1) {
      return value
  }
  const aiteIndex = value.indexOf('@')
  const pointIndex = value.lastIndexOf('.')
  const mail = value.substring(0, aiteIndex)
  if (mail.length <= 3) {
      result += '***'
  } else {
      result += value.substring(0, 3) + '***'
  }
  result += value.substring(pointIndex+1, value.length)
  return result;
}

const filters = {
  "dateFormat": dateFormat,
  "datetimeFormat": datetimeFormat,
  "errorFormat": errorFormat,
  "emailFormat": emailFormat,
};

export default {
  install(Vue) {
    Object.keys(filters).forEach(key => {
      Vue.filter(key, filters[key])
    });
  }
}
