package params

import (
	"fmt"
	"github.com/google/martian/log"
	"net/url"
)

var KubeApiServerUrlValidateFunc = func(p interface{}) error {
	val, ok := p.(string)
	if !ok {
		return fmt.Errorf(validateErrTypeParse, p, "string")
	}
	if val == "" {
		log.Debugf("kube-api-server not set,skip validate")
		return nil
	}
	_, err := url.Parse(val)
	if err != nil {
		return fmt.Errorf(validateUrlParse, val, err.Error())
	}
	return nil
}
