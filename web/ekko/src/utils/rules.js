import i18n from "@/i18n"

const RequiredRule = {
    required: true,
    trigger: "blur",
    message: i18n.t("commons.validate.required_msg"),
}
const NumberRule = {
    required: true,
    trigger: "blur",
    min: 1,
    type: "number",
    message: i18n.t("commons.validate.number_limit")
}
export default {
    RequiredRule,
    NumberRule,
}
