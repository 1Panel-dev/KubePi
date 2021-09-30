export function bystesLength (str) {
  let count = 0
  for (let i = 0; i < str.length; i++) {
    if (str.charCodeAt(i) > 255) {
      count += 2
    } else {
      count++
    }
  }
  return count
}

export function isJSON(str) {
  if (typeof str == 'string') {
    try {
      let obj=JSON.parse(str);
      return !!(typeof obj == 'object' && obj);

    } catch(e) {
      // console.log('error：'+str+'!!!'+e);
      return false;
    }
  }
  console.log('It is not a string!')
}
