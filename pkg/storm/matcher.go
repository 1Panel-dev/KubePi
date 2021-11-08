package storm

import (
	"github.com/asdine/storm/v3/q"
	"go/token"
	"reflect"
	"strings"
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

func Contains(fieldName string, val interface{}) q.Matcher {
	return q.NewFieldMatcher(fieldName, &contains{val: val})
}

type arrayValueEq struct {
	val string
}

func (c *arrayValueEq) MatchField(v interface{}) (bool, error) {
	refV := reflect.ValueOf(v)
	if refV.Kind() != reflect.Slice {
		return false, nil
	}
	valArray := strings.Split(c.val, ",")
	if len(valArray) != refV.Len() {
		return false, nil
	}

	var fieldArray []interface{}
	for i := 0; i < refV.Len(); i++ {
		fieldArray = append(fieldArray, refV.Index(i).Interface())
	}
	for i := range fieldArray {
		v, ok := fieldArray[i].(string)
		if !ok || v != valArray[i] {
			return false, nil
		}
	}
	return true, nil
}

func ArrayValueEq(fieldName, val string) q.Matcher {
	return q.NewFieldMatcher(fieldName, &arrayValueEq{val: val})
}

type like struct {
	val string
}

func (l *like) MatchField(v interface{}) (bool, error) {
	refV := reflect.ValueOf(v)
	if refV.Kind() == reflect.String {
		vs := v.(string)
		if strings.Contains(vs, l.val) {
			return true, nil
		}
	}
	return false, nil
}

func Like(fieldName string, val string) q.Matcher {
	return q.NewFieldMatcher(fieldName, &like{val: val})
}
