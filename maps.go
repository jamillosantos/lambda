package lambda

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type mapUtils map[string]string

func (m mapUtils) String(key string) (string, bool) {
	value, ok := m[key]
	return value, ok
}

func (m mapUtils) Int(key string) (int, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.Atoi(value)
}

func (m mapUtils) Int64(key string) (int64, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.ParseInt(value, 10, 64)
}

func (m mapUtils) Float64(key string) (float64, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.ParseFloat(value, 64)
}

func (m mapUtils) Bool(key string) (bool, error) {
	value, ok := m[key]
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return strconv.ParseBool(value)
}

type mapArrayUtils map[string][]string

func (m mapArrayUtils) Strings(key string) ([]string, error) {
	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return value, nil
}

func (m mapArrayUtils) Ints(key string) ([]int, error) {
	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result := make([]int, len(value))
	var err error
	for i, v := range value {
		result[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %d: %w", i, err)
		}
	}
	return result, nil
}

func (m mapArrayUtils) Int64s(key string) ([]int64, error) {
	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result := make([]int64, len(value))
	var err error
	for i, v := range value {
		result[i], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %d: %w", i, err)
		}
	}
	return result, nil
}

func (m mapArrayUtils) Float64s(key string) ([]float64, error) {
	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result := make([]float64, len(value))
	var err error
	for i, v := range value {
		result[i], err = strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %d: %w", i, err)
		}
	}
	return result, nil
}

func (m mapArrayUtils) Bools(key string) ([]bool, error) {
	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result := make([]bool, len(value))
	var err error
	for i, v := range value {
		result[i], err = strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %d: %w", i, err)
		}
	}
	return result, nil
}

type multiValues struct {
	mapUtils
	mapArrayUtils
}
