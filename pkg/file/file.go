package file

import (
	"k8s.io/client-go/util/homedir"
	"strings"
)

func ReplaceHomeDir(path string) string {
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", homedir.HomeDir(), -1)
	}
	return path
}
