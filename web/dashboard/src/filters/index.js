import {dateFormat, datetimeFormat} from "fit2cloud-ui/src/filters/time"

// export function DateFormat (value) {
//   let data = value
//   if (data == null) {
//     return null
//   }
//   let dt = new Date(data)
//   let yyyy = dt.getFullYear()
//   let MM = (dt.getMonth() + 1).toString().padStart(2, '0')
//   let dd = dt.getDate().toString().padStart(2, '0')
//   let h = dt.getHours().toString().padStart(2, '0')
//   let m = dt.getMinutes().toString().padStart(2, '0')
//   let s = dt.getSeconds().toString().padStart(2, '0')
//   return yyyy + '-' + MM + '-' + dd + ' ' + h + ':' + m + ':' + s
// }


export function AgeFormat (value) {
  if (value == null) {
    return ""
  }
  const begin = new Date(value)
  const end = new Date()
  const dateDiff = end.getTime() - begin.getTime()
  let dayDiff = Math.floor(dateDiff / (24 * 3600 * 1000)) //计算出相差天数
  let result = ""
  if (dayDiff > 0) {
    result = result + dayDiff + " days"
    return result
  }
  let leave1 = dateDiff % (24 * 3600 * 1000) //计算天数后剩余的毫秒数
  let hours = Math.floor(leave1 / (3600 * 1000)) //计算出小时数
  if (hours > 0) {
    result = result + hours + " hours"
    return result
  }
  let leave2 = leave1 % (3600 * 1000) //计算小时数后剩余的毫秒数
  let minutes = Math.floor(leave2 / (60 * 1000)) //计算相差分钟
  if (minutes > 0) {
    result = result + minutes + " mins"
    return result
  }
  //计算相差秒数
  let leave3 = leave2 % (60 * 1000); //计算分钟数后剩余的毫秒数
  let seconds = Math.round(leave3 / 1000);
  if (seconds > 0) {
    result = result + seconds + " secs"
    return result
  }
  return result
}

const filters = {
  "dateFormat": dateFormat,
  "datetimeFormat": datetimeFormat,
  "age": AgeFormat
}

export default {
  install (Vue) {
    Object.keys(filters).forEach(key => {
      Vue.filter(key, filters[key])
    })
  }
}
