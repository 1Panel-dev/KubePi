package i18n

var enUSMapping = TextMapping{
	"already exists":                        "Resource already exists. Try changing the resource name.",
	"username or password error":            "Login failed. The username or password is incorrect.",
	"Unauthorized":                          "Authorization error.",
	"permission %s required":                "Permission denied: %s",
	"please login":                          "Session has expired. Please log in.",
	"can not delete yourself":               "You cannot delete yourself.",
	"username can not be none":              "Username cannot be empty.",
	"must select one role":                  "You must select a role.",
	"must select one rule":                  "You must create a rule.",
	"user %s can not access resource %s %s": "User %s cannot access resource %s %s.",
	"can not match original password":       "The original password does not match.",
	"username already exists":               "Username already exists.",
	"email already exists":                  "Email already exists.",
	"unable to complete authorization":      "Unable to complete authorization. Check whether the username is valid: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/",
	"no login user":                         "No logged-in user.",
}
