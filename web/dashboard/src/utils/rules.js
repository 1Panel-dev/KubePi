import i18n from "@/i18n"

const RequiredRule = {
  required: true,
  trigger: "blur",
  message: i18n.t("commons.validate.required_input_msg"),
}
const SelectRule = {
  required: true,
  trigger: "change",
  message: i18n.t("commons.validate.required_select_msg"),
}
// 支持小写英文、数字和-
const CommonNameRule = {
  required: true,
  pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/,
  message: i18n.t("commons.validate.name_rules"),
  trigger: "blur"
}
// 支持数字或%
const NumberPersentRule = {
  required: true,
  pattern: /^[0-9]([0-9]*$|([0-9]*?%$))/,
  message: i18n.t("commons.validate.number_persent_rules"),
  trigger: "blur"
}
const NumberRule = {
  required: true,
  trigger: "blur",
  min: 1,
  type: "number",
  message: i18n.t("commons.validate.number_limit")
}
const UrlRule = {
  required: true,
  pattern: /^(https?):\/\/([a-zA-Z0-9.-]+(:[a-zA-Z0-9.&%$-]+)*@)*((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])){3}|([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.(com|edu|gov|int|mil|net|org|biz|arpa|info|name|pro|aero|coop|museum|[a-zA-Z]{2}))(:[0-9]+)*(\/($|[a-zA-Z0-9.,?'\\+&%$#=~_-]+))*$/,
  message: i18n.t("commons.validate.url_rules"),
  trigger: "blur"
}
export default {
  RequiredRule,
  SelectRule,
  NumberPersentRule,
  NumberRule,
  CommonNameRule,
  UrlRule,
}
