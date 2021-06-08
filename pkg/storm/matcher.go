package storm

import (
	"github.com/asdine/storm/v3/q"
	"go/token"
	"reflect"
)

type contains struct {
	val interface{}
}

func (c *contains) MatchField(v interface{}) (bool, error) {
	refV := reflect.ValueOf(v)
	if refV.Kind() != reflect.Slice {
		return false, nil
	}
	t := token.EQL
	for i := 0; i < refV.Len(); i++ {
		if compare(c.val, refV.Index(i).Interface(), t) {
			return true, nil
		}
	}
	return false, nil
}

func Containers(fieldName string, val interface{}) q.Matcher {
	return q.NewFieldMatcher(fieldName, &contains{val: val})
}
