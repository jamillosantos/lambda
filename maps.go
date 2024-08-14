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
	if !ok {
		return "", false
	}
	return value, true
}

func (m mapUtils) StringDefault(key, def string) string {
	value, ok := m.String(key)
	if !ok {
		return def
	}
	return value
}

func (m mapUtils) Int(key string) (int, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapUtils) IntDefault(key string, def int) int {
	value, err := m.Int(key)
	if err != nil {
		return def
	}
	return value
}

func (m mapUtils) Int64(key string) (int64, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapUtils) Int64Default(key string, def int64) int64 {
	value, err := m.Int64(key)
	if err != nil {
		return def
	}
	return value
}

func (m mapUtils) Float64(key string) (float64, error) {
	value, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapUtils) Float64Default(key string, def float64) float64 {
	value, err := m.Float64(key)
	if err != nil {
		return def
	}
	return value
}

func (m mapUtils) Bool(key string) (bool, error) {
	value, ok := m[key]
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s: %w", key, err)
	}
	return result, nil
}

func (m mapUtils) BoolDefault(key string, def bool) bool {
	value, err := m.Bool(key)
	if err != nil {
		return def
	}
	return value
}
