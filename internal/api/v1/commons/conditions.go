package commons

import "github.com/1Panel-dev/KubePi/internal/service/v1/common"

type SearchConditions struct {
	Conditions common.Conditions `json:"conditions"`
}
