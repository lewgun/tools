package provider

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	drivers = map[string]interface{}{}
	Default = &Console{}
)

func register(driver string, constructor interface{}) error {

	lower := strings.ToLower(driver)
	if _, isExists := drivers[lower]; isExists {
		return fmt.Errorf("The driver: %s is existed\n", driver)
	}

	drivers[lower] = constructor
	return nil
}
func New(typ string, arg *Arg) interface{} {

	i, _ := call(drivers, typ, arg)
	return i
}

func call(m map[string]interface{}, name string, params ...interface{}) (result interface{}, err error) {

	f := reflect.ValueOf(m[name])

	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}

	in := make([]reflect.Value, len(params))

	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)[0].Interface()
	return
}
