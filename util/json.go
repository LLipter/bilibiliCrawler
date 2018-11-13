package util

import (
	"errors"
	"fmt"
)

func MissingError(key string) error {
	return errors.New(fmt.Sprintf("missing '%s'", key))
}

func TypeError(key string) error {
	return errors.New(fmt.Sprintf("'%s' type error", key))
}

func JsonGetDict(dict map[string]interface{}, key string) (map[string]interface{}, error) {
	obj, ok := dict[key]
	if !ok {
		return nil, MissingError(key)
	}
	res, ok := obj.(map[string]interface{})
	if !ok {
		return nil, TypeError(key)
	}
	return res, nil
}

func JsonGetArray(dict map[string]interface{}, key string) ([]interface{}, error) {
	obj, ok := dict[key]
	if !ok {
		return nil, MissingError(key)
	}
	res, ok := obj.([]interface{})
	if !ok {
		return nil, TypeError(key)
	}
	return res, nil
}

func JsonGetStr(dict map[string]interface{}, key string) (string, error) {
	obj, ok := dict[key]
	if !ok {
		return "", MissingError(key)
	}
	res, ok := obj.(string)
	if !ok {
		return "", TypeError(key)
	}
	return res, nil
}

func JsonGetInt64(dict map[string]interface{}, key string) (int64, error) {
	obj, ok := dict[key]
	if !ok {
		return 0, MissingError(key)
	}
	res, ok := obj.(float64)
	if !ok {
		return 0, TypeError(key)
	}
	return int64(res), nil
}
