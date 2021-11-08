package storm

import (
	"github.com/asdine/storm/v3/q"
	"go/token"
	"reflect"
	"sort"
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

type arrayValueLike struct {
	val string
}

func (a *arrayValueLike) MatchField(v interface{}) (bool, error) {
	refV := reflect.ValueOf(v)
	if refV.Kind() != reflect.Slice {
		return false, nil
	}
	vs := strings.Split(a.val, ",")
	for i := range vs {
		var flag = false
		for j := 0; j < refV.Len(); j++ {
			v, ok := refV.Index(j).Interface().(string)
			if ok {
				if v == vs[i] {
					flag = true
				}
			}
		}
		if !flag {
			return false, nil
		}
	}
	return true, nil
}

func ArrayValueLike(fieldName, val string) q.Matcher {
	return q.NewFieldMatcher(fieldName, &arrayValueLike{val: val})
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

	var fieldArray []string
	for i := 0; i < refV.Len(); i++ {
		if refV.Index(i).Kind() == reflect.String {
			fieldArray = append(fieldArray, refV.Index(i).Interface().(string))
		} else {
			return false, nil
		}
	}
	sort.Strings(fieldArray)
	sort.Strings(valArray)

	for i := range fieldArray {
		if fieldArray[i] != valArray[i] {
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
