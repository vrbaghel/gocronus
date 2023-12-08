package env

import (
	"errors"
	"fmt"
	"reflect"
)

func checkENVErrors(errSlice []error) error {
	for _, err := range errSlice {
		if err != nil {
			return errors.New(fmt.Sprint("environment misconfigured; error", err.Error()))
		} else {
			continue
		}
	}
	return nil
}

func validateENVs(envs ...interface{}) error {
	for _, env := range envs {
		valuesSlice := reflect.ValueOf(env)
		for i := 0; i < valuesSlice.NumField(); i++ {
			if valuesSlice.Field(i).Kind() == reflect.String && valuesSlice.Field(i).Interface() == "" {
				return errors.New(fmt.Sprint("environment misconfigured; missing field", valuesSlice.Type().Field(i).Name))
			}
		}
	}
	return nil
}
