package util

import (
	"errors"
	"fmt"
)

func missingError(key string) error{
	return errors.New(fmt.Sprintf("missing '%s'", key))
}

func typeError(key string) error{
	return errors.New(fmt.Sprintf("'%s' type error", key))
}

func JsonGetDict(dict map[string]interface{}, key string) (map[string]interface{}, error){
	obj, ok := dict[key]
	if !ok {
		return nil, missingError(key)
	}
	res, ok := obj.(map[string]interface{})
	if !ok {
		return nil, typeError(key)
	}
	return res, nil
}

func JsonGetStr(dict map[string]interface{}, key string) (string, error){
	obj, ok := dict[key]
	if !ok {
		return "", missingError(key)
	}
	res, ok := obj.(string)
	if !ok {
		return "", typeError(key)
	}
	return res, nil
}

func JsonGetInt64(dict map[string]interface{}, key string) (int64, error){
	obj, ok := dict[key]
	if !ok {
		return 0, missingError(key)
	}
	res, ok := obj.(float64)
	if !ok {
		return 0, typeError(key)
	}
	return int64(res), nil
}