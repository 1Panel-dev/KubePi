import i18n from "@/i18n"

const RequiredRule = {
    required: true,
    trigger: "blur",
    message: i18n.t("commons.validate.required"),
}

const EmailRule = {
    type: 'email',
    message: i18n.t("commons.validate.email"),
    trigger: ['blur']
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

const LengthRule = {
    min: 1,
    max: 30,
    message: i18n.t("commons.validate.limit", [1, 30]),
    trigger: "blur"
}
const NamePattern = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/

const CommonNameRule = {
    required: true,
    pattern: NamePattern,
    message: i18n.t("commons.validate.name_not_compliant"),
    trigger: "blur"
}


export default {
    RequiredRule,
    NumberRule,
    EmailRule,
    PasswordRule,
    CommonNameRule,
    LengthRule
}
