package bean_util

import (
	"errors"
	"github.com/goinggo/mapstructure"
	"reflect"
	"strings"
)

func ProToStandard(in interface{}, out interface{}) (err error) {
	if in == nil || out == nil {
		err = errors.New("[ProToStandard] in or out is nil")
	}
	m := make(map[string]interface{})
	elem := reflect.ValueOf(&in).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		if !strings.Contains(relType.Field(i).Name, "XXX") {
			m[relType.Field(i).Name] = elem.Field(i).Interface()
		}
	}
	//将 map 转换为指定的结构体
	err = mapstructure.Decode(m, out)
	return
}
