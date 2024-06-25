// to unit `m`
export function cpuUnitConvert(value) {
  if (value.indexOf("u") !== -1) {
    return Number((Number(value.replace("u", "")) / 1000000000).toFixed(2))
  } else if (value.indexOf("n") !== -1) {
    return Number((Number(value.replace("n", "")) / 1000000).toFixed(2))
  } else if (value.indexOf("m") !== -1) {
    return Number(value.replace("m", ""))
  } else {
    return Number(value) * 1000
  }
}

// to unit `Mi`
export function memoryUnitConvert(value) {
  if (value.indexOf("Ki") !== -1) {
    return Number((Number(value.replace("Ki", "")) / 1024).toFixed(2))
  } else if (value.indexOf("Mi") !== -1) {
    return Number(value.replace("Mi", ""))
  } else if (value.indexOf("Gi") !== -1) {
    return Number(value.replace("Gi", "")) * 1024
  } else if (!isNaN(parseFloat(value)) && isFinite(value)) {
    // 如果value是一个纯数字，则以字节为单位的，并转换为MiB
    return Number((value / 1024 / 1024).toFixed(2))
  } else if (value.indexOf("G") !== -1) {
    return Number(value.replace("G", "")) * 1000
  } else if (value.indexOf("m") !== -1) {
    return Number(value.replace("m", "")) / 1000.00 / 1024.00 /1024.00 
  } else if (value.indexOf("M") !== -1) {
    return Number(value.replace("M", ""))  * 1024.00 /1000.00
  } else {
    return Number(value)
  }
}


