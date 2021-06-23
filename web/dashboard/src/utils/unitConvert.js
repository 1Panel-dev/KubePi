export function cpuUnitConvert(value) {
  if (value.indexOf("m") !== -1) {
    return value.replace("m", "")
  } else {
    return Number(value) * 1000
  }
}

export function memeryUnitConvert(value) {
  if (value.indexOf("Gi") !== -1) {
    return Number(value.replace("Gi", "")) * 1024
  } else if (value.indexOf("Mi") !== -1) {
    return Number(value.replace("Mi", ""))
  } else {
    return Number(value)
  }
}