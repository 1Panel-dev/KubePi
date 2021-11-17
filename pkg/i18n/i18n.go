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
		return "", fmt.Errorf("can not translate to language %s, key: %s", lang, key)
	}
	msg := mapping[key]
	if msg == "" {
		return "", fmt.Errorf("can not find i18n key %s", key)
	}
	if len(v) > 0 {
		ks := v[0].([]string)

		var iks []interface{}
		for i := range ks {
			iks = append(iks, ks[i])
		}
		return fmt.Sprintf(msg, iks...), nil
	} else {
		return msg, nil
	}

}
