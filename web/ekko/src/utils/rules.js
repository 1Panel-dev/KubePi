import i18n from "@/i18n"

const RequiredRule = {
    required: true,
    trigger: "blur",
    message: i18n.t("commons.validate.required"),
}

const EmailRule = {
    type: 'email',
    message: i18n.t("commons.validate.email"),
    trigger: ['blur', 'change']
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
    EmailRule
}
