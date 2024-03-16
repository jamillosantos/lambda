package lambda

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type mapArrayUtils map[string][]string

func (m mapArrayUtils) Strings(key string) ([]string, error) {
	value, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	return value, nil
}

func (m mapArrayUtils) String(key string) (string, bool) {
	value, ok := m[key]
	if !ok {
		return "", false
	}
	if len(value) == 0 {
		return "", false
	}
	return value[len(value)-1], true
}

func (m mapArrayUtils) StringDefault(key, def string) string {
	value, ok := m.String(key)
	if !ok {
		return def
	}
	return value
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

func (m mapArrayUtils) Int(key string) (int, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	if len(value) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.Atoi(value[len(value)-1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapArrayUtils) IntDefault(key string, def int) int {
	value, err := m.Int(key)
	if err != nil {
		return def
	}
	return value
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

func (m mapArrayUtils) Int64(key string) (int64, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	if len(value) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.ParseInt(value[len(value)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapArrayUtils) Int64Default(key string, def int64) int64 {
	value, err := m.Int64(key)
	if err != nil {
		return def
	}
	return value
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

func (m mapArrayUtils) Float64(key string) (float64, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	if len(value) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.ParseFloat(value[len(value)-1], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapArrayUtils) Float64Default(key string, def float64) float64 {
	value, err := m.Float64(key)
	if err != nil {
		return def
	}
	return value
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

func (m mapArrayUtils) Bool(key string) (bool, error) {
	value, ok := m[key]
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	if len(value) == 0 {
		return false, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.ParseBool(value[len(value)-1])
	if err != nil {
		return false, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapArrayUtils) BoolDefault(key string, def bool) bool {
	value, err := m.Bool(key)
	if err != nil {
		return def
	}
	return value
}
