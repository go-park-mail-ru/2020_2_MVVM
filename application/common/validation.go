package common

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"reflect"
)

const (
	validateErr = "validation failed with error code: "
	nameTag     = "alphanum~1,The name should contain alphanumeric characters.,stringlength(4|15)~2,The name must contain between 4 and 15 characters."
)


//TODO: если использование санитайзера действительно оправдано, определить все разрешенные тэги и символы, добавить обработку строчных слайсов
func clearHtml(req interface{}) {
	var clearText string

	p := bluemonday.UGCPolicy()
	p.AllowStandardURLs()
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	length := val.NumField()
	for i := 0; i < length; i++ {
		valField := val.Field(i)
		if valField.Kind() == reflect.String {
			valStr := valField.String()
			if valStr != "" {
				clearText = p.Sanitize(valStr)
				if clearText != valStr {
					valField.SetString(clearText)
				}
			}
		}
	}
}

func ReqValidation(req interface{},) error {
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return errors.New(validateErr + err.Error())
	}
	//clearHtml(req)
	return nil
}
