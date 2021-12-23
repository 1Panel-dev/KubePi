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
  } else {
    return Number(value)
  }
}

