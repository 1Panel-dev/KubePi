package i18n

import "fmt"

const (
	LanguageZhCN = "zh-CN"
	LanguageEnUS = "en-US"
)

type TextMapping map[string]string

type LangMap map[string]TextMapping

var langMapping = map[string]TextMapping{
	LanguageZhCN: zhCNMapping,
	LanguageEnUS: enUSMapping,
}

func Translate(lang string, key string, v ...interface{}) (string, error) {
	mapping := langMapping[lang]
	if mapping == nil {
		return "", fmt.Errorf("can no trspate to language %s", lang)
	}
	msg := mapping[key]
	if msg == "" {
		return "", fmt.Errorf("can not find i18n key %s", key)
	}
	return fmt.Sprintf(msg, v...), nil
}
