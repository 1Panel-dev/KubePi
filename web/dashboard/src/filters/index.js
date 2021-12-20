import {dateFormat, datetimeFormat} from "fit2cloud-ui/src/filters/time"

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

export function CpuFormat (value) {
  if (value == null) {
    return ""
  }
  if (value.indexOf("u") !== -1) {
    return (Number(value.replace("u", "")) / 1000000000).toFixed(2) + "m"
  } else if (value.indexOf("n") !== -1) {
    return (Number(value.replace("n", "")) / 1000000).toFixed(2) + "m"
  } else if (value.indexOf("m") !== -1) {
    return value
  } else {
    return (Number(value) * 1000) + "m"
  }
}

export function MemoryFormat (value) {
  if (value == null) {
    return ""
  }
  if (value.indexOf("Ki") !== -1) {
    return (Number(value.replace("Ki", "")) / 1024).toFixed(2) + "Mi"
  } else if (value.indexOf("Mi") !== -1) {
    return value
  } else if (value.indexOf("Gi") !== -1) {
    return (Number(value) * 1000) + "Mi"
  }
}

const filters = {
  "dateFormat": dateFormat,
  "datetimeFormat": datetimeFormat,
  "age": AgeFormat,
  "cpu": CpuFormat,
  "memory": MemoryFormat
}

export default {
  install (Vue) {
    Object.keys(filters).forEach(key => {
      Vue.filter(key, filters[key])
    })
  }
}
