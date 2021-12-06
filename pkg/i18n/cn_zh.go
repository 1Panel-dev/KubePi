package i18n

var zhCNMapping = TextMapping{
	"already exists":                        "资源已存在,请尝试修改资源名称",
	"username or password error":            "登录失败,用户名或密码错误",
	"Unauthorized":                          "认证失败",
	"permission %s required":                "权限不被允许:%s",
	"please login":                          "会话失效，请重新登录",
	"can not delete yourself":               "无法删除您自己",
	"username can not be none":              "用户名不能为空",
	"must select one role":                  "请至少选择一个角色",
	"must select one rule":                  "请至少创建一个规则",
	"user %s can not access resource %s %s": "用户 %s 缺少资源 [%s - %s] 的权限, 无法完成此操作",
	"can not match original password":       "无法匹配原密码",
	"username already exists":               "用户名已存在",
	"email already exists":                  "邮箱已存在",
	"unable to complete authorization":      "无法完成授权，请检查用户名是否符合规范: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/",
}
