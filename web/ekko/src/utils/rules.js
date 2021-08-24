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

const PasswordPattern = /^(?=.*\d)(?=.*[a-zA-Z])[\da-zA-Z~!@#$%^&*]{8,30}$/

const PasswordRule = {
    required: true,
    pattern: PasswordPattern,
    message: i18n.t("commons.validate.password_help"),
    trigger: "blur"
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
    EmailRule,
    PasswordRule
}
